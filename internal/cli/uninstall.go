package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Force bool

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&Force, "force", "f", false, "Force removal of installed packages")
}

var removeCmd = &cobra.Command{
	Use:   "uninstall [flags] packages...",
	Short: "Removes the installed configs",
	Long:  "Removes the config files for installed programs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		removeInstalled(args)
		return nil
	},
}

func removeInstalled(packages []string) {
	for i := 0; i < len(packages); i++ {
		sharedir := utils.ShareDir() // Directory of file
		pkg := packages[i]           // Current package
		pkgrcloc := path.Join(sharedir, fmt.Sprintf("%s.yaml", pkg))
		pkgrcfile, err := os.ReadFile(pkgrcloc)
		utils.DieIfError(err, fmt.Sprintf("%s is not installed or may have been installed incorrectly or with an older version of the Catppuccin CLI. Please reinstall and try again.", pkg))
		var pkgrc structs.AppLocation
		err = yaml.Unmarshal(pkgrcfile, &pkgrc)
		utils.DieIfError(err, fmt.Sprintf("Failed to read saved data for %s. Error: %s", pkg, err))
		remove := pkgrc.Location
		for e := 0; e < len(remove); e++ {
			fmt.Printf("Removing %s...\n", remove[e])
			os.Remove(remove[e])
		}
		os.Remove(pkgrcloc) // Remove the pkgrc
		fmt.Println("Finished!")
	}
}
