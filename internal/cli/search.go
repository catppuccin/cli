package cli

import (
	"fmt"
	"os"

	// "github.com/catppuccin/cli/internal/pkg/structs"
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
		fmt.Printf("%s\n", searchQuery[i])
	}
	dir := utils.ShareDir() + "/repos.json"
	if !utils.PathExists(dir) { // If repos.json file does not exist.
		fmt.Printf("\n%s does not exist. Caching JSON...", dir)
		utils.UpdateJSON()
	}
	body, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println("Cannot open file. ")
		os.Exit(1)
	}
	for i := 0; i < len(searchQuery); i++ {
		cache, err := structs.UnmarshalSearch(body)
		if err != nil {
			fmt.Printf("Error opening cache: %s", err)
		} else {
			result := utils.SearchRepos(cache, searchQuery[i])
			var resp string
			fmt.Printf("Found repo: %s", result.Name)
			fmt.Println("Do you want to install it? (Y/n)")
			if fmt.Scanln(&resp); resp == "Y" || resp == "y" {
				installer(searchQuery)
			}
		}
	}
}
