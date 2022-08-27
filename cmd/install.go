package cmd

import (
	"fmt"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"runtime"
)

var Flavour string
var Modes string

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&Flavour, "flavour", "f", "all", "Custom flavour")
	installCmd.Flags().StringVarP(&Modes, "mode", "m", "default", "Custom mode")
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
			//continue
			os.Exit(1)
		} else {
			success = append(success, string(repo))
		}
	}

	fmt.Println("\nChecking for installed packages:")
	programs := []structs.Program{}
	programLocations := []string{}
	programNames := []string{}

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
			fmt.Println(ctprc.AppName)
			InstallDir := ""
			if runtime.GOOS == "windows" {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Windows)
				fmt.Printf(InstallDir)
			} else if runtime.GOOS == "linux" {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Linux)
				fmt.Printf(InstallDir)
			} else {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Macos) // Just make the naive assumption that if it's not Windows or Linux, it's MacOS.
				fmt.Printf(InstallDir)
			}
			_, err = os.Stat(InstallDir)

			if err != nil {
				fmt.Printf("%s was not detected. %s\n", InstallDir, err)
			} else {
				fmt.Printf("%s path found at %s", ctprc.AppName, InstallDir)

				programs = append(programs, ctprc)
				programLocations = append(programLocations, InstallDir)
				programNames = append(programNames, success[i])
			}
		}
	}
	for i := 0; i < len(programs); i++ {
		fmt.Println("\nCloning " + programs[i].AppName + "...")
		baseDir := utils.CloneRepo(programNames[i])
		ctprc := programs[i]
		//Symlink the repo
		switch Flavour {
		// TO-DO: Implement modes
		case "all":
			utils.MakeLinks(baseDir, ctprc.Installation.InstallFlavours.All.Default, ctprc.Installation.To, programLocations[i]) // The magic line
		case "latte":
			utils.MakeLinks(baseDir, ctprc.Installation.InstallFlavours.Latte.Default, ctprc.Installation.To, programLocations[i])
		case "frappe":
			utils.MakeLinks(baseDir, ctprc.Installation.InstallFlavours.Frappe.Default, ctprc.Installation.To, programLocations[i])
		case "macchiato":
			utils.MakeLinks(baseDir, ctprc.Installation.InstallFlavours.Macchiato.Default, ctprc.Installation.To, programLocations[i])
		case "mocha":
			utils.MakeLinks(baseDir, ctprc.Installation.InstallFlavours.Mocha.Default, ctprc.Installation.To, programLocations[i])
		}
	}
}
