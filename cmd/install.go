package cmd

import (
	"catppuccin/uwu/internal/pkg/structs"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use: "install", 
	Short: "Install programs",
	Long: `Installs the programs listed from the official Catppuccin repos.`,
	Run: func(cmd *cobra.Command, args []string) {
		installer(args)
	},
}

func installer(packages []string) {

	fmt.Println("Installing the follow packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}

	fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccinrc
		rc := fmt.Sprintf("https://raw.githubusercontent.com/catppuccin/%s/main/.ctprc", repo)
		res, err := http.Get(rc)
		if err != nil {
			fmt.Printf("\nFailed to make HTTP request: %s\n", err)
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("%s does not have a .ctprc.\n", repo)
			continue
		} else {
			success = append(success, string(repo))
		}
	}

	fmt.Println("\nChecking for installed packages...")
	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/catppuccin/" + success[i] + "/main/.ctprc"
		res, err := http.Get(rc)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(res.Body)

		body, err := io.ReadAll(res.Body)

		programs := []structs.Program{}

		if err != nil {
			fmt.Println("Failed to read body.")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			if err != nil {
				fmt.Printf("Failed to parse .ctprc for %s\n", success[i])
			}
			fmt.Println(ctprc)
			path, err := exec.LookPath(ctprc.PathName)
			if err != nil {
				// Program is not installed/could not be detected
				fmt.Printf("%s was not detected.\n", ctprc.PathName)
			} else {
				fmt.Printf("%s found at location %s.\n", ctprc.PathName, path)
				programs = append(programs, ctprc)
			}
		}
	}
}

func genChezmoi(repo string, dir string, refresh int) string {
	// Creates a chezmoi entry using the repo name, updates every week
	res := fmt.Sprintf("%s:\n  type: git-repo\n  url: \"%s.git\"\n  refreshPeriod: %dh\n\n", dir, repo, refresh)
	return res
}
