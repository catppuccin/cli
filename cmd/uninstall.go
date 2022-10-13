package cmd

import (
	"fmt"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"runtime"
)

var Force bool

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&Force, "force", "f", false, "Force removal of installed packages")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the installed configs",
	Long:  "Removes the config files for installed programs",
	Run: func(cmd *cobra.Command, args []string) {
		removeInstalled(args)
	},
}

func removeInstalled(packages []string) {
	fmt.Println("Checking if the packages are installed...")
	for i := 0; i < len(packages); i++ {
		fmt.Printf("%s\n", packages[i])
	}
	for i := 0; i < len(packages); i++ {
		stageDir := path.Join(utils.ShareDir(), packages[i]) // stage directory
		fmt.Println(stageDir)
		ctpYaml := stageDir + "/.catppuccin.yaml" // Set the directory for .catppuccin.yaml and read it.
		yamlContent, err := os.ReadFile(ctpYaml)
		if err != nil {
			fmt.Println("\nCould not read file.")
			os.Exit(1)
		}
		ctprc, err := structs.UnmarshalProgram(yamlContent)
		if err != nil {
			fmt.Println("\nCould not unmarshal file.")
		}
		fileLoc := "" // Determine the OS and set the installation direction from .catppuccin.yaml
		if runtime.GOOS == "windows" {
			fileLoc = utils.HandleDir(ctprc.Installation.InstallLocation.Windows)
		} else if runtime.GOOS == "linux" {
			fileLoc = utils.HandleDir(ctprc.Installation.InstallLocation.Linux)
		} else {
			fileLoc = utils.HandleDir(ctprc.Installation.InstallLocation.Macos)
		}
		finalDir := path.Join(fileLoc, ctprc.Installation.To) // Sets the final direction for installation
		// repoDir := path.Join(utils.ShareDir(), packages[i])
		_, err = os.Lstat(finalDir)                           // Check for existence of the file and remove it.
		if err != nil {
			fmt.Printf("Could not find %s.", finalDir)
		} else {
			fmt.Printf("\nFound %s! Removing...", finalDir)
			err := os.RemoveAll(finalDir)
			if err != nil {
				fmt.Printf("Could not remove %s.", finalDir)
			}
			if Force == true { // If Force is set to true, we remove the staging directory for the file too.
				fmt.Printf("\nFound %s...", stageDir)
				err := os.RemoveAll(stageDir)
				if err != nil {
					fmt.Printf("Could not remove %s.", stageDir)
				}
			}
		}
	}
}
