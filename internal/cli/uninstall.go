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
	rootCmd.AddCommand(RemoveCmd)
	RemoveCmd.Flags().BoolVarP(&Force, "Force", "F", false, "Force removal of installed packages")
}

var RemoveCmd = &cobra.Command{
	Use:   "uninstall [flags] packages...",
	Short: "Removes the installed configs",
	Long:  "Removes the config files for installed programs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		RemoveInstalled(args)
		return nil
	},
}

func RemoveInstalled(packages []string) {
	for i := 0; i < len(packages); i++ {
		sharedir := utils.ShareDir() // Directory of file
		pkg := packages[i]           // Current package
		pkgrcloc := path.Join(sharedir, fmt.Sprintf("%s.yaml", pkg))
		ctprcloc := path.Join(sharedir, fmt.Sprintf("%s/.catppuccin.yaml", pkg))
		pkgrcfile, err := os.ReadFile(pkgrcloc)
		ctprcfile, err := os.ReadFile(ctprcloc)
		if err != nil {
			log.Fatalf("Could not read %s.yaml", pkg)
		}
		var pkgrc structs.AppLocation
		var ctprc structs.Program
		err = yaml.Unmarshal(pkgrcfile, &pkgrc)
		if err != nil {
			log.Fatalf("Failed to read saved data for %v.", pkg)
		}
		err = yaml.Unmarshal(ctprcfile, &ctprc)
		if err != nil {
			log.Fatalf("Failed to read .catppuccin.yaml data for %v.", pkg)
		}
		log.Infof("Running uninstall hooks...")
		hooks := ctprc.Installation.Hooks.Post.Uninstall
		for x := 0; x < len(hooks); x++ {
			if err := hooks[x].Run(); err != nil {
				log.Fatalf("Failed to run hook.")
				log.Debugf("%s", err)
			}
		}
		remove := pkgrc.Location
		for e := 0; e < len(remove); e++ {
			log.Infof("Removing %s...", remove[e])
			os.Remove(remove[e])

		}
		if Force == true {
			RepoLoc := path.Join(sharedir, pkg)
			log.Info("Removing cloned directory!")
			os.RemoveAll(RepoLoc)
			log.Info("Deleted cloned Repo!")
		}
		os.Remove(pkgrcloc) // Remove the pkgrc
		log.Info("Finished!")
	}
}
