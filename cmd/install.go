package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var Flavour string

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&Flavour, "flavour", "f", "all", "Custom flavour")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the config",
	Long:  "Installs the configs by cloning them from the Catppuccin repos.",
	Run: func(cmd *cobra.Command, args []string) {
		installer(args)
	},
}

func installer(packages []string) {
	// Hard-coded variables, will add functionality to change these via flags
	//mode := "default"
	fmt.Println("Installing packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	//fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccin.yaml
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.catppuccin.yaml", org, repo)
		fmt.Printf(rc)
		res, err := http.Get(rc)
		if err != nil {
			fmt.Println("\nCould not make GET request")
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("\n%s does not have a .catppuccin.yaml", repo)
			continue
		} else {
			success = append(success, string(repo))
		}
	}

	fmt.Println("\nChecking for installed packages:")
	programs := []structs.Program{}
	programsLocations := []string{}
	programsURLs := []string{}

	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/" + org + "/" + success[i] + "/main/.catppuccin.yaml"
		res, err := http.Get(rc)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println()
			}
		}(res.Body)

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Failed to read body")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			if err != nil {
				fmt.Printf(".catppuccin.yaml couldn't be unmarshaled correctly. Some data may be corrupted. (%s)\n", err)
			}
			fmt.Println(ctprc.Properties.AppName)
			InstallDir := ""
			if runtime.GOOS == "windows" {
				InstallDir = handleDir(ctprc.Properties.InstallLocation.Windows)
				fmt.Printf(InstallDir)
			} else {
				InstallDir = handleDir(ctprc.Properties.InstallLocation.Unix) // Just make the naive assumption that if it's not Windows, it's Unix.
				fmt.Printf(InstallDir)
			}
			_, err = os.Stat(InstallDir)

			if err != nil {
				fmt.Printf("%s was not detected. %s\n", InstallDir, err)
			} else {
				fmt.Printf("%s path found at %s", ctprc.Properties.AppName, InstallDir)

				programs = append(programs, ctprc)
				programsLocations = append(programsLocations, InstallDir)
				programsURLs = append(programsURLs, success[i])

			}
		}
	}
	for i := 0; i < len(programs); i++ {
		fmt.Println("\nCloning " + programs[i].Properties.AppName + "...")
		baseDir := cloneRepo(programs[i].Properties.AppName)
		ctprc := programs[i]
		//Symlink the repo
		switch Flavour {
		// TO-DO: Implement modes
		case "all":
			makeLinks(baseDir, ctprc.Installation.InstallFlavours.All.Default, ctprc.Installation.To, programsLocations[i]) // The magic line
		case "latte":
			makeLinks(baseDir, ctprc.Installation.InstallFlavours.Latte.Default, ctprc.Installation.To, programsLocations[i])
		case "frappe":
			makeLinks(baseDir, ctprc.Installation.InstallFlavours.Frappe.Default, ctprc.Installation.To, programsLocations[i])
		case "macchiato":
			makeLinks(baseDir, ctprc.Installation.InstallFlavours.Macchiato.Default, ctprc.Installation.To, programsLocations[i])
		case "mocha":
			makeLinks(baseDir, ctprc.Installation.InstallFlavours.Mocha.Default, ctprc.Installation.To, programsLocations[i])
		}
	}
}

func handleDir(dir string) string {
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

func makeLinks(baseDir string, links []string, to string, finalDir string) {
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
			handleFilePath(finalDir, name)
		} else {
			handleDirPath(finalDir, name)
		}
		fmt.Printf("Linking: %s to %s via %s\n", link, finalDir, name)
		// Use the name as name, the link as the from, and the finalDir as the to
		makeLink(link, finalDir, name)
	}
}

func handleDirPath(finalDir string, name string) {
	// Check if dir to link already exists
	fullDir := path.Join(finalDir, name)
	var resp string
	if utils.PathExists(fullDir) {
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

func handleFilePath(finalDir string, name string) {
	// Check if dir to link already exists
	fileFolder, _ := path.Split(name)
	fullDir := path.Join(finalDir, fileFolder)
	if !utils.PathExists(fullDir) {
		err := os.Mkdir(fullDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create parent directory %s", fullDir)
		}
	}
}

func cloneRepo(repo string) string {
	stagePath := path.Join(shareDir(), repo)
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      "https://github.com/catppuccin-rfc/" + repo + ".git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
	}
	return stagePath
}

func shareDir() string {
	if utils.IsWindows() {
		return path.Join(UserHomeDir(), "AppData/LocalLow/catppuccin-cli")
	}
	return path.Join(utils.GetEnv("XDG_DATA_HOME", handleDir("~/.local/")), "share/catppuccin-cli")
}

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
