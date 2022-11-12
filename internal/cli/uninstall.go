package cli

import (
	"fmt"
	"path"

	//"github.com/catppuccin/cli/internal/pkg/structs"
	"os"

	"github.com/caarlos0/log"
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
	Use:   "remove",
	Short: "Removes the installed configs",
	Long:  "Removes the config files for installed programs",
	Run: func(cmd *cobra.Command, args []string) {
		removeInstalled(args)
	},
}

func removeInstalled(packages []string) {
	for i := 0; i < len(packages); i++ {
		sharedir := utils.ShareDir() // Directory of file
		pkg := packages[i]           // Current package
		pkgrcloc := path.Join(sharedir, fmt.Sprintf("%s.yaml", pkg))
		pkgrcfile, err := os.ReadFile(pkgrcloc)
		if err != nil {
			log.Fatalf("Could not read %s.yaml", pkg)
		}
		var pkgrc structs.AppLocation
		err = yaml.Unmarshal(pkgrcfile, &pkgrc)
		if err != nil {
			log.Fatalf("Failed to read saved data for %v.", pkg)
		}
		remove := pkgrc.Location
		for e := 0; e < len(remove); e++ {
			log.Infof("Removing %s...", remove[e])
			os.Remove(remove[e])
		}
		os.Remove(pkgrcloc) // Remove the pkgrc
		log.Info("Finished!")
	}
}
