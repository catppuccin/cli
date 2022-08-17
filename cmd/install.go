package cmd

import (
	"fmt"
	"path"
	"runtime"
	"github.com/go-git/go-git/v5"
	"io"
	"net/http"
	"catppuccin/uwu/internal/utils"
	"os"
	"os/user"
	"strings"
	"catppuccin/uwu/internal/pkg/structs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
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
	fmt.Println("Installing packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Printf(packages[i])
	}
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccinrc
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.ctprc", org, repo)
		res, err := http.Get(rc)
		if err != nil {
			fmt.Println("\nCould not make GET request")
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("%s does not have a .ctprc", repo)
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
		rc := "https://raw.githubusercontent.com/" + org + "/" + success[i] + "/main/.ctprc"
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
				fmt.Println(".ctprc couldn't be unmarshaled correctly. Some data may be corrupted.")
			}
			fmt.Println(ctprc)
			InstallDir := ""
			if runtime.GOOS == "windows" {
				InstallDir = handleDir(ctprc.InstallLocation.Windows)
			} else {
				InstallDir = handleDir(ctprc.InstallLocation.Unix) // Just make the naive assumption that if it's not Windows, it's Unix.
			}
			_, err = os.Stat(InstallDir)

			if err != nil {
				fmt.Printf("%s was not detected. %s\n", InstallDir, err)
			} else {
				fmt.Printf("%s path found at %s", ctprc.AppName, InstallDir)

				programs = append(programs, ctprc)
				programsLocations = append(programsLocations, InstallDir)
				programsURLs = append(programsURLs, success[i])

				for i := 0; i < len(programs); i++ {
					createStagingDir(programs[i].PathName)
					
					fmt.Println("\nCloning " + programs[i].AppName + "...")

					cloneRepo(programs[i].AppName)
				}
			}
		}
	}
}

func createStagingDir(repo string) {
	err := os.MkdirAll(path.Join("stage/", repo), 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Failed to create staging directory")
		os.Exit(1)
	}
}

func handleDir(dir string) string {
	usr, _ := user.Current()
	dir = strings.Replace(dir, "~", usr.HomeDir, -1)
	appdata, _ := os.UserConfigDir()
	dir = strings.Replace(dir, "%appdata%", appdata, -1)
	return dir
}

func cloneRepo(repo string) {
	stagePath := path.Join(shareDir(), repo)
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      "https://github.com/catppuccin/" + repo + ".git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func shareDir() string {
	if utils.IsWindows() {
		return path.Join(UserHomeDir(), "AppData/LocalLow/uwu")
	}
	return path.Join(utils.GetEnv("XDG_DATA_HOME", handleDir("~/.local/")), "share/uwu")
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
