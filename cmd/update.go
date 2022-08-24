package cmd

import (
	"fmt"
	"github.com/catppuccin/cli/internal/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the local JSON cache",
	Long:  "Update the local JSON cache used to check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		updateJSON()
	},
}

func updateJSON() {
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	uc := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	req, err := http.Get(uc)
	if err != nil {
		fmt.Println("Could not make get request")
		os.Exit(1)
	}
	if req.StatusCode != 200 {
		fmt.Println("Could not get response")
		os.Exit(1)

	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println()
			}
		}(req.Body)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("\nCould not read response body")
			os.Exit(1)
		}
		dir := shareDir() + "/repos.json"   // Set the staging directory plus the file name
		err = os.WriteFile(dir, body, 0644) // Write the file.
		if err != nil {
			fmt.Println("Could not save the file. ")
		}
	}
}
