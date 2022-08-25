package cmd

import (
	"fmt"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/spf13/cobra"
	"os"
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
	dir := shareDir() + "/repos.json"
	if _, err := os.Stat(dir); err != nil { // If repos.json file does not exist.
		fmt.Printf("\n%s does not exist. Please run `ctp update` first.\n", dir)
		os.Exit(1)
	}
	fileJSON, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println("Cannot open file. ")
		os.Exit(1)
	}
	//qr := []structs.SearchRes{}
	//qc, err := structs.UnmarshalSearch(fileJSON)
	//qc := structs.SearchRes{}
	//var qc structs.SearchRes
	fmt.Println(fileJSON)
	qc, err := structs.UnmarshalSearch(fileJSON) // Please figure this part out I have picked my brains on this for some hours now. :sadcat:
	if err != nil {
		fmt.Printf("Cannot unmarshal file. %s", err)
		os.Exit(1)
	}
	for i := 0; i < len(qc.Name); i++ {
		fmt.Printf("%s\n", qc.Name)
	}
	for i := 0; i < len(searchQuery); i++ {
		if qc.Name == searchQuery[i] {
			var resp string
			if fmt.Scan(&resp); resp == "y" {
				//installer([]string{searchQuery[i]}) // We'll pass the installer function here.
			}
		} else {
			fmt.Println("Could not find the package.")
		}
	}
}
