package main

import (
	"fmt"
	"os"
	"net/http"
)

func main() {
	packages := os.Args[1:]

	fmt.Println("Installing the follow packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}
	fmt.Println("Generating chezmoi config...")
	success := []string{}
	for i := 0; i < len(packages); i++ {
		repo := packages[i]
		// Attempt to get the .catppuccinrc
		rc := "https://raw.githubusercontent.com/catppuccin/" + repo + "/main/.catppuccinrc"
		res, err := http.Get(rc)
		if err != nil {
			fmt.Printf("Failed to make HTTP request: %s\n", err)
			os.Exit(1)
		}
		if res.StatusCode != 200 {
			fmt.Printf("%s does not have a .catppuccinrc.\n", repo)
			break
		}
		success = append(success, repo)
	}
	fmt.Println(success)
}

func gen_chezmoi(repo string, dir string, refresh int) string {
	// Creates a chezmoi entry using the repo name, updates every week
	res := fmt.Sprintf("%s:\n  type: git-repo\n  url: \"%s.git\"\n  refreshPeriod: %dh\n", dir, repo, refresh)
	return res
}
