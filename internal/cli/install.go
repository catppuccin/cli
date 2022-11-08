package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	Flavour string
	Mode    string
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&Flavour, "flavour", "f", "all", "Custom flavour")
	installCmd.Flags().StringVarP(&Mode, "mode", "m", "default", "Custom mode")
}

var installCmd = &cobra.Command{
	Use:   "install [flags] packages...",
	Short: "Installs the config",
	Long:  "Installs the configs by cloning them from the Catppuccin repos.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		installer(args)
		return nil
	},
}

func installer(packages []string) {
	// Hard-coded variables, will add functionality to change these via flags
	fmt.Println("Installing packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	// fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccin.yaml
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.catppuccin.yaml", org, repo)
		fmt.Println(rc)
		res, err := http.Get(rc)
		if err != nil {
			fmt.Println("\nCould not make GET request")
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("%s does not have a .catppuccin.yaml\n", repo)
			// continue
			os.Exit(1)
		} else {
			success = append(success, string(repo))
		}
	}

	fmt.Println("\nChecking for installed packages:")
	programs := []structs.Program{}
	programLocations := []string{}
	programNames := []string{}
	comments := []string{}

	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/" + org + "/" + success[i] + "/main/.catppuccin.yaml"
		res, err := http.Get(rc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Failed to read body")
			continue
		}
		ctprc, err := structs.UnmarshalProgram(body)
		if err != nil {
			fmt.Printf(".catppuccin.yaml couldn't be unmarshaled correctly. Some data may be corrupted. (%s)\n", err)
			continue
		}
		fmt.Println(ctprc.AppName)
		var installDir string
		switch runtime.GOOS {
		case "windows":
			installDir = utils.HandleDir(ctprc.Installation.InstallLocation.Windows)
		case "linux":
			installDir = utils.HandleDir(ctprc.Installation.InstallLocation.Linux)
		case "darwin":
			installDir = utils.HandleDir(ctprc.Installation.InstallLocation.Macos)
		default:
			fmt.Println("current OS is not supported")
			continue
		}
		fmt.Println(installDir)
		_, err = os.Stat(installDir)
		if err != nil {
			fmt.Printf("%s was not detected. %s\n", installDir, err)
			continue
		}
		fmt.Printf("%s path found at %s\n", ctprc.AppName, installDir)
		programs = append(programs, ctprc)
		programLocations = append(programLocations, installDir)
		programNames = append(programNames, success[i])
		comments = append(comments, ctprc.Installation.Comments)

	}
	for i := 0; i < len(programNames); i++ {
		fmt.Println("\nCloning " + programs[i].AppName + "...")
		programName := programNames[i]
		installDir := path.Join(utils.ShareDir(), programName)
		baseDir := utils.CloneRepo(installDir, programName)
		installLoc := programLocations[i]
		ctprc := programs[i]
		// Symlink the repo
		returnedLocation := utils.InstallFlavours(baseDir, Mode, Flavour, ctprc, installLoc)
		utils.MakeLocation(packages[i], returnedLocation)
		if comments[i] != "" {
			fmt.Println("\nNote: " + comments[i])
		}
		for _, hook := range ctprc.Installation.Hooks.Install {
			if err := hook.Run(); err != nil {
				fmt.Println(err)
			}
		}
	}
	// nya~
}
