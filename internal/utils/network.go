/*
NETWORK.GO
Contains bigger functions
that have to do with Git
or networking.
*/
package utils

import ( // {{{
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	// "path/filepath"

	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v47/github"
) // }}}

// PullUpdates opens a git repo and pulls the latest changes.
func PullUpdates(repo string) { // {{{
	// Repo should be a valid folder, so now we'll open the .git
	r, err := git.PlainOpen(repo) // Open new repo targeting the .git folder
	if err != nil {
		fmt.Printf("Error opening repo folder: %s\n", err)
	} else {
		// Get working directory
		w, err := r.Worktree()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\n", err)
		} else {
			// Pull the latest changes from origin
			fmt.Printf("Pulling latest changes for %s...\n", repo)
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				fmt.Printf("Failed to pull updates: %s\n", err)
			}
		}
	}
} // }}}

// UpdateJSON makes a search request for all Catppuccin repos and caches them.
func UpdateJSON() { // {{{
	dir := path.Join(ShareDir(), "repos.json") // Set the staging directory plus the file name
	org := GetEnv("ORG_OVERRIDE", "catppuccin")
	client := github.NewClient(nil)

	// Get all the Catppuccin repositories
	opt := &github.RepositoryListByOrgOptions{Type: "public"} // Get all the repositories
	repos, _, err := client.Repositories.ListByOrg(context.Background(), org, opt)

	// Handle errors
	if err != nil {
		fmt.Println("Failed to get repositories.")
	} else {
		fmt.Println("Received repositories. Caching!")
		themes := []structs.SearchEntry{}
		for i := 0; i < len(repos); i++ {
			repo := repos[i]
			if !ListContains(repo.Topics, "catppuccin-meta") { // Repo does not contain catppuccin-meta topic
				// Append search result
				theme := structs.SearchEntry{
					Name:   repo.GetName(),
					Stars:  repo.GetStargazersCount(),
					Topics: repo.Topics,
				}
				themes = append(themes, theme)
			}
		}
		body, err := json.Marshal(themes)
		if err != nil {
			fmt.Printf("Failed to marshal cache: %s\nPlease try again.\n", err)
		} else {
			os.WriteFile(dir, body, 0o644)
		}
	}
} // }}}
