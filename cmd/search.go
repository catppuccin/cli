package cmd

import (
	"fmt"
	"os"

// "github.com/catppuccin/cli/internal/pkg/structs"
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
	_, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println("Cannot open file. ")
		os.Exit(1)
	}
	//qr := []structs.SearchRes{}
	//qc, err := structs.UnmarshalSearch(fileJSON)
	//qc := structs.SearchRes{}
	//var qc structs.SearchRes
	utils.UpdateJSON()
}
