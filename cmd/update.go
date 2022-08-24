package cmd

import (
	"fmt"
	"path"

	"github.com/catppuccin/cli/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// Create the command
func init() {
	rootCmd.AddCommand(updateCmd)
}


var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages.",
	Long:  "Checks repos for any updates to packages.",
	Run: func(cmd *cobra.Command, args []string) {
		updater(args)
	},
}

func updater(packages []string) {
	// Handle updates
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		repo_dir := path.Join(utils.ShareDir(), repo)
		if !utils.PathExists(repo_dir) {
			fmt.Printf("Cannot find installed %s.\n", repo)
		} else {
			utils.PullUpdates(repo_dir)
		}
	}
}

