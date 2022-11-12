package cli

import (
	"fmt"
	"os"

	// "github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for repositories for a package",
	Long:  "Searches for repositories for a package and provides with the option to install the theme",
	Run: func(cmd *cobra.Command, args []string) {
		searchPackage(args)
	},
}

func searchPackage(searchQuery []string) {
	fmt.Println("Searching for packages:")
	for i := 0; i < len(searchQuery); i++ {
		log.Infof("\n", searchQuery[i])
	}
	dir := utils.ShareDir() + "/repos.json"
	if !utils.PathExists(dir) { // If repos.json file does not exist.
		log.Infof("\n%s does not exist. Caching JSON...", dir)
		utils.UpdateJSON()
	}
	body, err := os.ReadFile(dir)
	if err != nil {
		log.WithError(err).Fatalf("Cannot open file. ")
	}
	for i := 0; i < len(searchQuery); i++ {
		cache, err := structs.UnmarshalSearch(body)
		if err != nil {
			log.WithError(err).Fatalf("Error opening cache: %s", err)
		} else {
			result := utils.SearchRepos(cache, searchQuery[i])
			var resp string
			log.Infof("Found repo: %s", result.Name)
			log.Info("\nDo you want to install it? (Y/n)")
			if fmt.Scanln(&resp); resp == "Y" || resp == "y" {
				installer(searchQuery)
			}
		}
	}
}
