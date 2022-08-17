package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
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
	//fmt.Printf("\nGenerating chezmoi config\n")

	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		rc := fmt.Sprintf("https://raw.githubusercontent.com/catppuccin/%s/main/.ctprc", repo)
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
		rc := "https://raw.githubusercontent.com/catppuccin/" + success[i] + "/main/.ctprc"
		res, _ := http.Get(rc)
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

			usr, _ := user.Current()
			ctprc, _ := structs.UnmarshalProgram(body)
			ctprc.InstallLocation = strings.Replace(ctprc.InstallLocation, "~", usr.HomeDir, -1)
			_, err := os.Stat(ctprc.InstallLocation)

			if err != nil {
				fmt.Printf("%s was not detected. %s\n", ctprc.InstallLocation, err)
			} else {
				fmt.Printf("%s path found at %s", ctprc.AppName, ctprc.InstallLocation)

				programs = append(programs, ctprc)
				programsLocations = append(programsLocations, ctprc.InstallLocation)
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

func cloneRepo(repo string) {
	dir, _ := os.Getwd()
	stagePath := path.Join(dir, "stage", repo)
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      "https://github.com/catppuccin/" + repo + ".git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
	}
}
