package cmd

import (
	"catppuccin/uwu/internal/pkg/structs"
	"fmt"
	"path"
	"io"
	"log"
	"net/http"
	"catppuccin/uwu/internal/utils"
	"os"
	"os/exec"
	"strings"

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
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccinrc
		rc := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.ctprc", org, repo)
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
	programs := []structs.Program{}
	programLocations := []string{}
	programNames := []string{}
	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/" + org + "/" + success[i] + "/main/.ctprc"
		res, err := http.Get(rc)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(res.Body)

		body, err := io.ReadAll(res.Body)


		if err != nil {
			fmt.Println("Failed to read body.")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			fmt.Println(ctprc)
			path, err := exec.LookPath(ctprc.PathName)
			if err != nil {
				// Program is not installed/could not be detected
				fmt.Printf("%s was not detected.\n", ctprc.PathName)
			} else {
				fmt.Printf("%s found at location %s.\n", ctprc.PathName, path)
				// Append program to detected programs and add it's location to a seperate list.
				programs = append(programs, ctprc)
				programLocations = append(programLocations, path)
				programNames = append(programNames, success[i])
			}
		}
	}
	// Part 3, clone the repo into staging dir
	for i := 0; i < len(programs); i++ {
		//loc := programLocations[i]
		ctprc := programs[i]
		programName := programNames[i]
		//installLoc := handleDir(loc, ctprc.InstallLocation)
		fmt.Println(fmt.Sprint(ctprc.InstallFiles))
		err := createStagingDir(programNames[i]) // Create directory with repo name
		if err != nil {
			// Directory already exists
			fmt.Printf("Directory for %s already exists! Skipping %s...\n(Run uwu update to update packages.)", programName, programName)
		} 
	}  
}

func createStagingDir(repo string) error {
	err := os.MkdirAll(path.Join("stage", repo), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	} else if err != nil {
		return err
	}
	// hey spoopy :)
	return nil
}

func handleDir(fileLoc string, dir string) string {
	// If install location begins with ./, replace the . with the fileLoc
	if strings.HasPrefix(dir, "./") {
		//Remove the dot, replace with fileLoc
		dir = fmt.Sprintf("%s%s", fileLoc, dir[1:]) 
	}
	if dir == "." {
		dir = fileLoc
	}
	return dir
}

func wrapQuotes(items []string) []string {
	for i := 0; i < len(items); i++ {
		items[i] = fmt.Sprintf("\"%s\"", items[i])
	}
	return items
}


