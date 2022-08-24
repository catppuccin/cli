package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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
	if strings.Contains(dir, "%userprofile%") { // For programs which store config on per-user basis like vscode
		dir = strings.Replace(dir, "%userprofile", usr.HomeDir, -1)
		fmt.Printf(dir)
	}
	dir = strings.Replace(dir, "%userprofile", usr.HomeDir, -1)
	dir = strings.Replace(dir, "~", usr.HomeDir, -1)
	appdata, _ := os.UserConfigDir()
	dir = strings.Replace(dir, "%appdata%", appdata, -1)
	return dir

}

// MakeLink makes a symlink from a path to another path with a suffix.
func makeLink(from string, to string, name string) {
	if to[len(to)-1:] != "/" {
		fmt.Println("'to' is not a directory wtf")
	} else {
		// Symlink the directory
		err := os.Symlink(from, path.Join(to, name)) /* Example:
		 * (Folder)
		 * Symlink themes/default into ~/.config/helix/themes
		 * from: ~/.local/share/catppuccin-cli/Helix/themes/default
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
}

// MakeLinks loops through a list and converts it's attributes into arguments for MakeLink.
func MakeLinks(baseDir string, links []string, to string, finalDir string) {
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
	for i := 0; i < len(links); i++ {
		link := path.Join(baseDir, links[i])
		// Use the regex to get the last part of the file URL and append it to the `to`
		shortPath := re.FindString(link)
		name := to
		if strings.Contains(shortPath[2:], ".") {
			// Path is a file, handle that
			name = path.Join(to, shortPath)
			HandleFilePath(finalDir, name)
		} else {
			HandleDirPath(finalDir, name)
		}
		fmt.Printf("Linking: %s to %s via %s\n", link, finalDir, name)
		// Use the name as name, the link as the from, and the finalDir as the to
		makeLink(link, finalDir, name)
	}
}

// HandleDirPath is a function to handle a directory when making a symlink
func HandleDirPath(finalDir string, name string) {
	// Check if dir to link already exists
	fullDir := path.Join(finalDir, name)
	var resp string
	if PathExists(fullDir) {
		fmt.Printf("Directory %s already exists.\nWould you like to move the directory?(y/N): ", fullDir)
		if fmt.Scan(&resp); resp == "y" {
			fmt.Println("\nReplacing directory...")
			prefix, suffix := path.Split(fullDir)
			renamed := suffix + "-" + time.Now().Format("06-01-02")
			renamed = path.Join(prefix, renamed)
			err := os.Rename(fullDir, renamed)
			if err != nil {
				fmt.Println("Failed to move directory. You may have to rerun this command with elevated permissions, or the old directory may already exist.")
				fmt.Printf("(Error: %s)\n", err)
			}
		}
	}
}

// HandleFilePath handles files that are made with symlinks
func HandleFilePath(finalDir string, name string) {
	// Check if dir to link already exists
	fileFolder, _ := path.Split(name)
	fullDir := path.Join(finalDir, fileFolder)
	if !PathExists(fullDir) {
		err := os.Mkdir(fullDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create parent directory %s", fullDir)
		}
	}
}

// CloneRepo clones a repo into the specified location.
func CloneRepo(repo string) string {
	stagePath := path.Join(ShareDir(), repo)
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      "https://github.com/" + GetEnv("ORG_OVERRIDE", "catppuccin") +  repo + ".git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
	}
	return stagePath
}

// PullUpdates opens a git repo and pulls the latest changes.
func PullUpdates(repo string) {
	// Repo should be a valid folder, so now we'll open the .git
	r, err := git.PlainOpen(repo)	 // Open new repo targeting the .git folder
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
