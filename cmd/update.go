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
			pullUpdates(repo_dir)
		}
	}
}

func pullUpdates(repo string) {
	// Repo should be a valid folder, so now we'll open the .git
	r, err := git.PlainOpen(repo)	 // Open new repo targeting the .git folder
	if err != nil {
		fmt.Printf("Error opening repo folder: %s\n", err)
	} else {
		// Get working directory
		w, err := r.Worktree()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\n", err)
		} else {
			// Pull the latest changes from origin
			fmt.Printf("Pulling latest changes for %s...\n", repo)
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				fmt.Printf("Failed to pull updates: %s\n", err)
			}
		}
	}
}
