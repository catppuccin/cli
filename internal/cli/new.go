package cli

import (
	"github.com/catppuccin/cli/internal/ui"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Initializes a new project.",
	Long:  "Uses the Catppuccin template repos to interactively create a new theme.",
	Run: func(cmd *cobra.Command, args []string) {
		handleArgs(args)
	},
}

func handleArgs(args []string) {
	if len(args) >= 1 {
		utils.InitTemplate(args[0], args[1], args[2], args[3], args[4])
	} else {
		ui.Run()
	}
}
