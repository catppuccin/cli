package cli

import (
	"fmt"
	"path"

	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
)

// Create the command
func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [flags] packages...",
	Short: "Update packages.",
	Long:  "Checks repos for any updates to packages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		updater(args)
		return nil
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
