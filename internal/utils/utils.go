/*
UTILS.GO ------------
Contains utility functions
and other small, useful
things
*/
package utils

import (

	// "log"

	"os"
	"path"

	// "path/filepath"
	"runtime"

	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// GetEnv gets an environment variable.
// If not defined, it gets the fallback.
func GetEnv(lookup string, fallback string) string {
	if res, ok := os.LookupEnv(lookup); ok {
		return res
	}
	return fallback
}

// IsWindows checks if OS is Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, exists := os.Stat(path)
	if os.IsNotExist(exists) {
		return false
	}
	return true
}

// UserHomeDir gets the user's home directory
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// DieIfError kills the program if err is not nil.
// func DieIfError(err error, message string) {
// 	if err != nil {
// 		// log.Fatalf("%s.\nFailed with error: %v", message, err)
// 		log.WithError(err).Fatalf("%s", message)
// 	}
// }

// ListContains checks if a list of strings contains a string
func ListContains(list []string, contains string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == contains {
			return true
		}
	}
	return false
}

// CheckBetter checks if better is greater than check. If it is, it returns better, otherwise it returns check. It also returns a BoolAnd of checkbetter and if better > check.
func CheckBetter(check int, better int, checkbetter bool) (int, bool) {
	if better > check {
		return better, BoolAnd(true, checkbetter)
	}
	return check, BoolAnd(false, checkbetter)
}

// BoolAnd uses booleans in an AND operator
func BoolAnd(first bool, second bool) bool {
	if first || second {
		return true
	}
	return false
}

// SearchRepos searches through a SearchRes for the best match
func SearchRepos(repos structs.SearchRes, term string) structs.SearchEntry {
	var best structs.SearchEntry
	bestScore := -1000
	for i := 0; i < len(repos); i++ {
		repo := repos[i]
		better := false
		rank := fuzzy.RankMatch(term, repo.Name)
		bestScore, better = CheckBetter(bestScore, rank, better) // Sets the new best score and also tells if if new term is better
		for e := 0; e < len(repo.Topics); e++ {
			topic := repo.Topics[e]
			rank = fuzzy.RankMatch(term, topic)
			bestScore, better = CheckBetter(bestScore, rank, better) // Basically what this does is goes and tells us the best match of the topic, and sets that score in bestScore.
			// If better is true, best becomes this repo. Just trust me on this. Just trust me on this.
		}
		if better {
			best = repo
		}
	}
	return best // Return the best match
}

// CloneTemplate creates the template directory and clones the template repo into it.
func CloneTemplate(repo string) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory.")
	}

	// Make project directory and clone
	installPath := path.Join(cwd, repo)
	err = os.Mkdir(installPath, 0o755)
	if err != nil {
		log.Fatalf("Failed to create project directory for %v.", repo)
	}
	CloneRepo(installPath, "template") // Clone the template repo into the installPath
}

// GetTemplateDir gets the location of the template directory.
func GetTemplateDir(repo string) string {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory.")
	}
	installPath := path.Join(cwd, repo)
	return installPath
}

// RunHooks runs a list of hooks.
func RunHooks(hooks []structs.Hook) {
	for _, hook := range hooks {
		if err := hook.Run(); err != nil {
			log.Fatalf("Failed to run hook.")
			log.Debugf("%s", err)
		}
	}
}
