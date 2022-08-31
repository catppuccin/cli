package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/catppuccin/cli/internal/ui"
	"github.com/catppuccin/cli/internal/utils"
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
}

func createRepoAlt() {
  fmt.Println("Creating new repo...")
  // Get current directory
  cwd, err := os.Getwd()
  utils.DieIfError(err, "Failed to get current directory.")
  // Get project name
  var repo string
  fmt.Print("Project name: ")
  fmt.Scan(&repo)
  fmt.Print("\n")
  // Make project directory and clone
  installPath := path.Join(cwd, repo)
  err = os.Mkdir(installPath, 0755)
  utils.DieIfError(err, fmt.Sprintf("Failed to make project directory for %s.", repo))
  utils.CloneRepo(installPath, "template") // Clone the template repo into the installPath
}
