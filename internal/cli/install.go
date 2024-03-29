package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	Flavour string
	Mode    string
	// force   bool
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&Flavour, "flavour", "f", "all", "Custom flavour")
	installCmd.Flags().StringVarP(&Mode, "mode", "m", "default", "Custom mode")
	installCmd.Flags().BoolVarP(&Force, "Force", "F", false, "Forcefully remove the cloned repository and install")
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
	log.DecreasePadding()
	log.Info("Installing packages...")
	log.IncreasePadding()
	for i := 0; i < len(packages); i++ {
		log.Infof(packages[i])
	}
	log.ResetPadding()
	log.Info("Downloading packages...")
	log.IncreasePadding()
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccin.yaml
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.catppuccin.yaml", org, repo)
		res, err := http.Get(rc)
		if err != nil {
			log.Fatalf("Could not make GET request")
		}
		if res.StatusCode != 200 {
			log.Errorf("%s does not have a .catppuccin.yaml.", repo)
		} else {
			success = append(success, string(repo))
		}
	}
	log.ResetPadding()
	log.IncreasePadding()
	log.Info("Checking for installed packages:")
	programs := []structs.Program{}
	programLocations := []string{}
	programNames := []string{}
	comments := []string{}

	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/" + org + "/" + success[i] + "/main/.catppuccin.yaml"
		res, err := http.Get(rc)
		if err != nil {
			log.Errorf("Could not make GET request")
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Errorf("Failed to read body")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			if err != nil {
				log.Errorf(".catppuccin.yaml couldn't be unmarshaled correctly. Some data may be corrupted.\n")
			}
			var InstallDir string
			switch runtime.GOOS {
			case "windows":
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Windows)
			case "linux":
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Linux)
			case "darwin":
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Macos)
			default:
				log.Errorf("Your OS is not supported.")
			}
			_, err = os.Stat(InstallDir)

			if err != nil {
				log.Errorf("%s was not detected. \n", InstallDir)
			} else {
				log.Infof("%s path found at %s", ctprc.AppName, InstallDir)
				programs = append(programs, ctprc)
				programLocations = append(programLocations, InstallDir)
				programNames = append(programNames, success[i])
				comments = append(comments, ctprc.Installation.Comments)
			}
		}
	}
	for i := 0; i < len(programNames); i++ {
		programName := programNames[i]
		installDir := path.Join(utils.ShareDir(), programName)
		repoExistsLoc := path.Join(utils.ShareDir(), programName)
		installLoc := programLocations[i]
		ctprc := programs[i]
		_, err := os.Stat(repoExistsLoc)
		if err == nil && !Force {
			log.Fatalf("Repository already exists. Please run the install with the -F=true flag again. ")
		}
		if Force {
			log.Infof("Removing existing repository directory: %v", repoExistsLoc)
			log.Info("Re-cloning the directory. ")
			RemoveInstalled([]string{programName})
		}
		log.ResetPadding()
		if ctprc.Installation.Hooks.Pre.Install != nil {
			log.Info("Running pre-install hooks...")
			log.IncreasePadding()
			utils.RunHooks(ctprc.Installation.Hooks.Pre.Install)
			log.ResetPadding()
		}
		log.Info("Cloning " + programs[i].AppName + "...")
		log.IncreasePadding()
		baseDir := utils.CloneRepo(installDir, programName)

		//Symlink the repo
		returnedLocation := utils.InstallFlavours(baseDir, Mode, Flavour, ctprc, installLoc)
		utils.MakeLocation(packages[i], returnedLocation)
		if comments[i] != "" {
			log.Info(comments[i])
		}
		log.ResetPadding()
		if ctprc.Installation.Hooks.Post.Install != nil {
			log.Info("Running post-install hooks...")
			utils.RunHooks(ctprc.Installation.Hooks.Post.Install)
		}
	}
	// nya~
}
