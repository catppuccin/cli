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

var Flavour string
var Mode string

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&Flavour, "flavour", "f", "all", "Custom flavour")
	installCmd.Flags().StringVarP(&Mode, "mode", "m", "default", "Custom mode")
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
	log.DecreasePadding()
	log.Info("Installing packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccin.yaml
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.catppuccin.yaml", org, repo)
		res, err := http.Get(rc)
		if err != nil {
			log.WithError(err).Fatalf("Could not make GET request")
		}
		if res.StatusCode != 200 {
			log.WithError(err).Fatalf("%s does not have a .catppuccin.yaml", repo)
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
		res, _ := http.Get(rc)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.WithError(err).Fatalf("Could not close body")
			}
		}(res.Body)

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.WithError(err).Fatalf("Failed to read body")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			if err != nil {
				log.WithError(err).Fatalf(".catppuccin.yaml couldn't be unmarshaled correctly. Some data may be corrupted. (%s)\n", err)
			}
			fmt.Println(ctprc.AppName)
			InstallDir := ""
			if runtime.GOOS == "windows" {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Windows)
			} else if runtime.GOOS == "linux" {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Linux)
			} else {
				InstallDir = utils.HandleDir(ctprc.Installation.InstallLocation.Macos) // Just make the naive assumption that if it's not Windows or Linux, it's MacOS.
			}
			_, err = os.Stat(InstallDir)

			if err != nil {
				log.WithError(err).Fatalf("%s was not detected. %s\n", InstallDir, err)
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
		log.Infof("Cloning " + programs[i].AppName + "...")
		programName := programNames[i]
		installDir := path.Join(utils.ShareDir(), programName)
		baseDir := utils.CloneRepo(installDir, programName)
		installLoc := programLocations[i]
		ctprc := programs[i]

		//Symlink the repo
		returnedLocation := utils.InstallFlavours(baseDir, Mode, Flavour, ctprc, installLoc)
		utils.MakeLocation(packages[i], returnedLocation)
		if comments[i] != "" {
			os.Stdout.WriteString("\nNote: " + comments[i])
		}
	}
	// nya~
}
