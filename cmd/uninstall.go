package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func init() {
	rootCmd.AddCommand(removeCmd)
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
	fmt.Println("Checking if the packages are installed.")
	for i := 0; i < len(packages); i++ {
		fmt.Printf("%s\n", packages[i])
	}
	stageDir := shareDir() // Get the staging directory
	for i := 0; i < len(packages); i++ {
		dir := path.Join(stageDir, packages[i]) // Join staging directory and the package
		finalDir := handleDir(dir)              // Handle "%APPDATA%" and "~"
		_, err := os.Lstat(finalDir)            // Check if the directory exists
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s does not exist", finalDir)
				os.Exit(0)
			}
		} else {
			fmt.Printf("Removing %s\n", finalDir)
			err := os.RemoveAll(finalDir) // Remove the directory
			if err != nil {
				fmt.Printf("\nCould not remove %s", finalDir)
			}
		}
	}
}
