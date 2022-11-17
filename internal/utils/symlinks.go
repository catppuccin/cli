/*
SYMLINKS.GO ------------
Contains bigger functions
related to creating and
handling symlinks.
*/
package utils

import (
	"os"
	"path"

	// "path/filepath"
	"regexp"

	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/pkg/structs"
)

// MakeLink makes a symlink from a path to another path with a suffix.
func makeLink(from string, to string, name string) string {
	symfile := path.Join(to, name)
	if to[len(to)-1:] != "/" {
		// fmt.Printf("\n'%s' is not a directory.", to)
		log.Fatalf("'%s' is not a directory.", to)
	} else if PathExists(to + name) {
		log.Info("Symlink already exists, removing and relinking...")
		err := os.RemoveAll(to + name)

		/* This remove a directory and replaces it with the new synlink in case it already exists.
		The reason to use RemoveAll was that Remove cannot delete a directory if it is not empty.
		On Unix, RemoveAll uses `rm -rf Dir/`*/

		if err != nil {
			// fmt.Printf("Failed to remove symlink. (Error: %s)\n", err)
			log.Fatalf("Failed to remove symlink.")
		}
		// Now symlinking again
		err = os.Symlink(from, symfile)
		if err != nil {
			log.Error("Failed to symlink.")
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
			// fmt.Println(err)
			log.Error("Failed to symlink.")
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
	// fmt.Println("Making symlinks....")
	log.Info("Making symlinks....")
	// Regex last-item match
	re, _ := regexp.Compile(`\/[^\/]*$`)
	// Iterate over links and use makeLink to make the links
	var symfiles []string
	for i := 0; i < len(links); i++ {
		link := path.Join(baseDir, links[i])
		linkInfo, err := os.Stat(link)
		if err != nil {
			log.Error("Failed to get info about file.")
		}
		// Check for a file extension; literally just looks for a "."
		shortPath := re.FindString(link)
		name := to
		if !linkInfo.IsDir() {
			// Path is a file, handle that
			name = path.Join(to, shortPath)
			HandleFilePath(finalDir, name)
			// Just link the file
			log.Infof("Linking: %s to %s via %s", link, finalDir, name)
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
	if err != nil {
		log.Fatalf("Failed to read directory while parsing for symlinks.")
	}
	for i := 0; i < len(files); i++ {
		files[i] = path.Join(link, files[i])
	}
	return files
}

// HandleFilePath handles files that are made with symlinks
func HandleFilePath(finalDir string, name string) {
	// Check if dir to link already exists
	fileFolder, _ := path.Split(name)
	fullDir := path.Join(finalDir, fileFolder)
	if !PathExists(fullDir) {
		err := os.Mkdir(fullDir, 0o755)
		if err != nil {
			log.Error("Failed to create parent directory.")
		}
	}
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
		log.Error("Mode does not exist.")
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
		log.Fatalf("Unexpected flavour: %s", flavour)
	}
	return nil
}
