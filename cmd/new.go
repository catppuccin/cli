package cmd

import (
	"log"

	"github.com/catppuccin/cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
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
		createRepo()
	},
}

func createRepo() {
	p := tea.NewProgram(ui.InitialModel())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	spinner := tea.NewProgram(ui.InitialModelSpinner()) // Spinner code. We can remove this for now.
	// Just here for demonstration purpose right now.
	if err := spinner.Start(); err != nil {
		log.Fatal(err)
	}
}
