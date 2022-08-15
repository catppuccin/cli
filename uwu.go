package main

import (
	"catppuccin/installer/internal/pkg/structs"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	packages := os.Args[1:]

	fmt.Println("Installing the follow packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}

	fmt.Println("\nGenerating chezmoi config...")
	var success []string
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccinrc
		rc := "https://raw.githubusercontent.com/catppuccin/" + repo + "/main/.ctprc"
		res, err := http.Get(rc)
		if err != nil {
			fmt.Printf("\nFailed to make HTTP request: %s\n", err)
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("%s does not have a .ctprc.\n", repo)
			continue
		} else {
			success = append(success, repo)
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

		if err != nil {
			fmt.Println("Failed to read body.")
		} else {
			ctprc, err := structs.UnmarshalProgram(body)
			if err != nil {
				fmt.Printf("Failed to parse .ctprc for %s\n", success[i])
			}
			fmt.Println(ctprc)
			path, err := exec.LookPath(success[i])
			if err != nil {
				// Program is not installed/could not be detected
				fmt.Printf("%s was not detected.\n", success[i])
				success = removeAt(success, i)
			} else {
				fmt.Printf("%s found at location %s.\n", success[i], path)
			}
		}
	}
}

func genChezmoi(repo string, dir string, refresh int) string {
	// Creates a chezmoi entry using the repo name, updates every week
	res := fmt.Sprintf("%s:\n  type: git-repo\n  url: \"%s.git\"\n  refreshPeriod: %dh\n", dir, repo, refresh)
	return res
}

func removeAt(list []string, index int) []string {
	return append(list[:index], list[index+1:]...)
}
