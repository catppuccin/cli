package main

import (
	"fmt"
	"github.com/kylelemons/go-gypsy"
	"io"
	"net/http"
	"os"
	"os/exec"
)

/*
type Program struct {
	app_name string
	path_name string
	install_location string
	install_files []string
	install_flavours struct {
		latte []struct {
			preference string
		}
		frappe []struct {
			mode string
		}
		macchiato []struct {
			preference string
		}
		mocha []struct {
			preference string
		}
	}
	one_flavour bool
	mode []string
}
*/

func main() {
	packages := os.Args[1:]

	fmt.Println("Installing the follow packages...")
	for i := 0; i < len(packages); i++ {
		fmt.Println(packages[i])
	}
	fmt.Println("\nGenerating chezmoi config...\n")
	success := []string{}
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
	fmt.Println("\nChecking for installed packages...\n")
	for i := 0; i < len(success); i++ {
		rc := "https://raw.githubusercontent.com/catppuccin/" + success[i] + "/main/.ctprc"
		res, err := http.Get(rc)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Failed to read body.")
		} else {
			config, err := yaml.Read(body)
			if err != nil {
				fmt.Printf("Failed to parse .ctprc for %s", success[i])
			}
			fmt.Println(ctprc)
			path, err := exec.LookPath(success[i])
			if err != nil {
				// Program is not installed/could not be detected
				fmt.Printf("%s was not detected.\n", success[i])
				success = remove_at(success, i)
			} else {
				fmt.Printf("%s found at location %s.\n", success[i], path)
			}
		}
	}
}

func gen_chezmoi(repo string, dir string, refresh int) string {
	// Creates a chezmoi entry using the repo name, updates every week
	res := fmt.Sprintf("%s:\n  type: git-repo\n  url: \"%s.git\"\n  refreshPeriod: %dh\n", dir, repo, refresh)
	return res
}

func remove_at(list []string, index int) []string {
	return append(list[:index], list[index+1:]...)
}
