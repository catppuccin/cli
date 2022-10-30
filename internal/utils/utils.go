package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	// "path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v47/github"
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

// ShareDir generates the share directory for the cli.
func ShareDir() string {
	if IsWindows() {
		return path.Join(UserHomeDir(), "AppData/LocalLow/catppuccin-cli")
	}
	return path.Join(GetEnv("XDG_DATA_HOME", HandleDir("~/.local/")), "share/catppuccin-cli")
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

// HandleDir handles a directory, replacing certain parts with known attributes.
func HandleDir(dir string) string {
	usr, _ := user.Current()
	envre, err := regexp.Compile(`\$([A-z0-9_\-]+)\/`) // Create the regex to detect environment variables
	DieIfError(err, "Failed to compile environment checking regex. Try running again.")
	dir = strings.Replace(dir, "%userprofile%", usr.HomeDir, -1)
	dir = strings.Replace(dir, "~", usr.HomeDir, -1)
	appdata, _ := os.UserConfigDir()
	dir = strings.Replace(dir, "%appdata%", appdata, -1)
	envar := string(envre.Find([]byte(dir)))
	if envar != "" {
		result := GetEnv(envar[1:len(envar)-1], "/////////////") // Intentionally screws up the program if failed
		if result[len(result)-1:] != "/" {
			result += "/"
		}
		dir = envre.ReplaceAllString(dir, result)
	}
	return dir
}

// MakeLink makes a symlink from a path to another path with a suffix.
func makeLink(from string, to string, name string) string {
	symfile := path.Join(to, name)
	if to[len(to)-1:] != "/" {
		fmt.Printf("\n'%s' is not a directory.", to)
		os.Exit(1)
	} else if PathExists(to + name) {
		fmt.Println("Symlink already exists, removing and relinking...")
		err := os.RemoveAll(to + name)

		/* This remove a directory and replaces it with the new synlink in case it already exists.
		The reason to use RemoveAll was that Remove cannot delete a directory if it is not empty.
		On Unix, RemoveAll uses `rm -rf Dir/`*/

		if err != nil {
			fmt.Printf("Failed to remove symlink. (Error: %s)\n", err)
			os.Exit(1)
		}
		// Now symlinking again
		err = os.Symlink(from, symfile)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// Symlink the directory
		err := os.Symlink(from, symfile) /* Example:
		 * (Folder)cin-cli/Helix/them
		 * Symlink themes/default into ~/.config/helix/themes
		 * from: ~/.local/share/catppuces/default
		 * to:   ~/.config/helix/
		 * name: themes/
		 * Creates a symlink from ~/.local/share/catppuccin-cli/Helix/themes to ~/.config/helix/themes
		 * (File)
		 * Symlink themes/default/catppuccin_mocha.toml into ~/.config/helix/themes
		 * from: ~/.local/share/catppuccin-cli/Helix/themes/default/catppuccin_mocha.toml
		 * to:   ~/.config/helix/
		 * name: themes/catppuccin_mocha.toml
		 * Creates a symlink from ~/.local/share/catppuccin-cli/Helix/themes/default/catppuccin_mocha.toml to ~/.config/helix/themes/catppuccin_mocha.toml
		 */
		if err != nil {
			fmt.Println(err)
		}
	}
	return symfile
}

// MakeLinks loops through a list and converts its attributes into arguments for MakeLink.
func MakeLinks(baseDir string, links []string, to string, finalDir string) []string {
	/* An explanation of these ambiguous names
	 * baseDir  - the directory in which the repo was staged, returned by cloneRepo
	 * links    - a list of files that we loop through to make links of
	 * to       - the location these were meant to be linked to, not including the actual path
	 * finalDir - the actual path they are going to
	 */
	fmt.Println("Making symlinks....")
	// Regex last-item match
	re, _ := regexp.Compile(`\/[^\/]*$`)
	// Iterate over links and use makeLink to make the links
	var symfiles []string
	for i := 0; i < len(links); i++ {
		link := path.Join(baseDir, links[i])
		linkInfo, _ := os.Stat(link)
		// Check for a file extension; literally just looks for a "."
		shortPath := re.FindString(link)
		name := to
		if !linkInfo.IsDir() {
			// Path is a file, handle that
			name = path.Join(to, shortPath)
			HandleFilePath(finalDir, name)
			// Just link the file
			fmt.Printf("Linking: %s to %s via %s\n", link, finalDir, name)
			// Use the name as name, the link as the from, and the finalDir as the to
			symfiles = append(symfiles, makeLink(link, finalDir, name))
		} else {
			files := HandleDirPath(baseDir, links[i], finalDir, name)
			symfiles = MakeLinks(baseDir, files, to, finalDir)
		}
	}
	return symfiles
}

// HandleDirPath is a function to handle a directory when making a symlink
func HandleDirPath(baseDir string, link string, finalDir string, name string) []string {
	// The link directory
	from := path.Join(baseDir, link)
	files, err := OSReadDir(from)
	DieIfError(err, "Failed to read directory while parsing for symlinking.")
	for i := 0; i < len(files); i++ {
		files[i] = path.Join(link, files[i])
	}
	return files
}

func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

// HandleFilePath handles files that are made with symlinks
func HandleFilePath(finalDir string, name string) {
	// Check if dir to link already exists
	fileFolder, _ := path.Split(name)
	fullDir := path.Join(finalDir, fileFolder)
	if !PathExists(fullDir) {
		err := os.Mkdir(fullDir, 0o755)
		if err != nil {
			fmt.Printf("Failed to create parent directory %s", fullDir)
		}
	}
}

// CloneRepo clones a repo into the specified location.
func CloneRepo(stagePath string, repo string) string {
	org := GetEnv("ORG_OVERRIDE", "catppuccin")
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL: fmt.Sprintf("https://github.com/%s/%s.git", org, repo),
	})
	if err != nil {
		fmt.Println(err)
	}
	return stagePath
}

// DieIfError kills the program if err is not nil.
func DieIfError(err error, message string) {
	if err != nil {
		log.Fatalf("%s.\nFailed with error: %v", message, err)
	}
}

// PullUpdates opens a git repo and pulls the latest changes.
func PullUpdates(repo string) {
	// Repo should be a valid folder, so now we'll open the .git
	r, err := git.PlainOpen(repo) // Open new repo targeting the .git folder
	if err != nil {
		fmt.Printf("Error opening repo folder: %s\n", err)
	} else {
		// Get working directory
		w, err := r.Worktree()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\n", err)
		} else {
			// Pull the latest changes from origin
			fmt.Printf("Pulling latest changes for %s...\n", repo)
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				fmt.Printf("Failed to pull updates: %s\n", err)
			}
		}
	}
}

// ListContains checks if a list of strings contains a string
func ListContains(list []string, contains string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == contains {
			return true
		}
	}
	return false
}

// UpdateJSON makes a search request for all Catppuccin repos and caches them.
func UpdateJSON() {
	dir := path.Join(ShareDir(), "repos.json") // Set the staging directory plus the file name
	org := GetEnv("ORG_OVERRIDE", "catppuccin")
	client := github.NewClient(nil)

	// Get all the Catppuccin repositories
	opt := &github.RepositoryListByOrgOptions{Type: "public"} // Get all the repositories
	repos, _, err := client.Repositories.ListByOrg(context.Background(), org, opt)

	// Handle errors
	if err != nil {
		fmt.Println("Failed to get repositories.")
	} else {
		fmt.Println("Received repositories. Caching!")
		themes := []structs.SearchEntry{}
		for i := 0; i < len(repos); i++ {
			repo := repos[i]
			if !ListContains(repo.Topics, "catppuccin-meta") { // Repo does not contain catppuccin-meta topic
				// Append search result
				theme := structs.SearchEntry{
					Name:   repo.GetName(),
					Stars:  repo.GetStargazersCount(),
					Topics: repo.Topics,
				}
				themes = append(themes, theme)
			}
		}
		body, err := json.Marshal(themes)
		if err != nil {
			fmt.Printf("Failed to marshal cache: %s\nPlease try again.\n", err)
		} else {
			os.WriteFile(dir, body, 0o644)
		}
	}
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

// InstallLinks is a wrapper over MakeLinks that parses the mode and uses it to create the correct link, as specified by the ctprc.
func InstallLinks(baseDir string, entry structs.Entry, to string, finalDir string, mode string) []string {
	if mode == "default" {
		// Default mode, just run makeLinks
		return MakeLinks(baseDir, entry.Default, to, finalDir) // The magic line
	}
	// Mode code
	modes := entry.Additional
	modeEntry := modes[mode]
	if modeEntry == nil {
		fmt.Printf("Mode '%s' does not exist.\n", mode)
	} else {
		return MakeLinks(baseDir, modeEntry, to, finalDir)
	}
	return nil
}

// InstallFlavours is a wrapper for InstallLinks which takes the flavour and handles the install accordingly
func InstallFlavours(baseDir string, mode string, flavour string, ctprc structs.Program, installLoc string) []string {
	switch flavour {
	case "all":
		return InstallLinks(baseDir, ctprc.Installation.InstallFlavours.All, ctprc.Installation.To, installLoc, mode)
	case "latte":
		return InstallLinks(baseDir, ctprc.Installation.InstallFlavours.Latte, ctprc.Installation.To, installLoc, mode)
	case "frappe":
		return InstallLinks(baseDir, ctprc.Installation.InstallFlavours.Frappe, ctprc.Installation.To, installLoc, mode)
	case "macchiato":
		return InstallLinks(baseDir, ctprc.Installation.InstallFlavours.Macchiato, ctprc.Installation.To, installLoc, mode)
	case "mocha":
		return InstallLinks(baseDir, ctprc.Installation.InstallFlavours.Mocha, ctprc.Installation.To, installLoc, mode)
	default:
		log.Fatal("Unexpected flavour")
	}
	return nil
}

// CloneTemplate creates the template directory and clones the template repo into it.
func CloneTemplate(repo string) {
	// Get current directory
	cwd, err := os.Getwd()
	DieIfError(err, "Failed to get current directory.")

	// Make project directory and clone
	installPath := path.Join(cwd, repo)
	err = os.Mkdir(installPath, 0o755)
	DieIfError(err, fmt.Sprintf("Failed to make project directory for %s", repo))
	CloneRepo(installPath, "template") // Clone the template repo into the installPath
}

// GetTemplateDir gets the location of the template directory.
func GetTemplateDir(repo string) string {
	// Get current directory
	cwd, err := os.Getwd()
	DieIfError(err, "Failed to get current directory.")
	installPath := path.Join(cwd, repo)
	return installPath
}

// InitTemplate initializes a template repo for the repo name specified.
func InitTemplate(repo string, exec string, linuxloc string, macloc string, windowsloc string) {
	installPath := GetTemplateDir(repo)
	ctprc, err := os.OpenFile(path.Join(installPath, ".catppuccin.yaml"), os.O_WRONLY, 0o644)
	DieIfError(err, "Failed to open .catppuccin.yaml.")
	defer ctprc.Close()
	content, err := os.ReadFile(path.Join(installPath, ".catppuccin.yaml")) // Don't use ioutil.ReadFile. Deprecated.
	DieIfError(err, "Failed to read .catppuccin.yaml.")

	ctp, err := template.New("catppuccin").Parse(string(content))
	DieIfError(err, "Failed to parse .catppuccin.yaml.")
	catppuccin := structs.Catppuccinyaml{
		Name:          repo,
		Exec:          exec,
		MacosLocation: macloc,
		LinuxLocation: linuxloc,
		WinLocation:   windowsloc,
	}

	err = ctp.Execute(ctprc, catppuccin)
	DieIfError(err, fmt.Sprintf("Failed to write to .catppuccin.yaml:%s", err))
}

// MakeLocation saves the locations written to during installation into a file for later access.
func MakeLocation(packages string, location []string) {
	flavourrc := structs.AppLocation{
		Location: location,
	}
	marshallData, err := flavourrc.MarshalLocation()
	if err != nil {
		fmt.Println("Failed to marshall data.")
	}
	filepath := packages + ".yaml"
	finalPath := path.Join(ShareDir(), filepath)

	if PathExists(finalPath) { // If it already exists, remove it
		os.Remove(finalPath)
	}

	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Cannot open file.")
	}
	if _, err := file.Write(marshallData); err != nil {
		fmt.Println("Failed to write to file.")
	}
}
