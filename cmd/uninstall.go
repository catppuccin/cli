package cmd

import (
	"fmt"
	"path"
	//"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
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
		pkgrc := path.Join(sharedir, fmt.Sprintf("%s.yaml", pkg))
		fmt.Println(pkgrc)
	}
}
