/*
SYSTEMS.GO
Contains bigger functions
that have to deal with
paths or the filesystem.
*/
package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/caarlos0/log"
	"github.com/catppuccin/cli/internal/pkg/structs"
	"github.com/go-git/go-git/v5"
)

// HandleDir handles a directory, replacing certain parts with known attributes.
func HandleDir(dir string) string {
	usr, _ := user.Current()
	envre, err := regexp.Compile(`\$([A-z0-9_\-]+)\/`) // Create the regex to detect environment variables
	if err != nil {
		log.WithError(err).Fatalf("Failed to compile regex.")
	}
	dir = strings.Replace(dir, "%userprofile%", usr.HomeDir, -1)
	dir = strings.Replace(dir, "~", usr.HomeDir, -1)
	appdata, _ := os.UserConfigDir()
	dir = strings.Replace(dir, "%appdata%", appdata, -1)
	envar := string(envre.Find([]byte(dir)))
	if envar != "" {
		result := GetEnv(envar[1:len(envar)-1], "/////////////") // Intentionally screws up the program if failed
		if result[len(result)-1:] != "/" {
			result += "/"
		}
		dir = envre.ReplaceAllString(dir, result)
	}
	return dir
}

// ShareDir generates the share directory for the cli.
func ShareDir() string {
	if IsWindows() {
		return path.Join(UserHomeDir(), "AppData/LocalLow/catppuccin-cli")
	}
	return path.Join(GetEnv("XDG_DATA_HOME", HandleDir("~/.local/")), "share/catppuccin-cli")
}

// OSReadDir expands a directory path into a list of files within that directory. Not recursive.
func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

// CloneRepo clones a repo into the specified location.
func CloneRepo(stagePath string, repo string) string {
	org := GetEnv("ORG_OVERRIDE", "catppuccin")
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:               fmt.Sprintf("https://github.com/%s/%s.git", org, repo),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Error("Could not clone repo")
	}
	return stagePath
}

// InitTemplate initializes a template repo for the repo name specified.
func InitTemplate(repo string, exec string, linuxloc string, macloc string, windowsloc string) {
	installPath := GetTemplateDir(repo)
	ctprc, err := os.OpenFile(path.Join(installPath, ".catppuccin.yaml"), os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Failed to open .catppuccin.yaml")
	}
	defer ctprc.Close()
	content, err := os.ReadFile(path.Join(installPath, ".catppuccin.yaml")) // Don't use ioutil.ReadFile. Deprecated.
	if err != nil {
		log.Fatalf("Failed to read .catppuccin.yaml")
	}
	ctp, err := template.New("catppuccin").Parse(string(content))
	if err != nil {
		log.Fatalf("Failed to parse .catppuccin.yaml")
	}
	catppuccin := structs.Catppuccinyaml{
		Name:          repo,
		Exec:          exec,
		MacosLocation: macloc,
		LinuxLocation: linuxloc,
		WinLocation:   windowsloc,
	}

	err = ctp.Execute(ctprc, catppuccin)
	if err != nil {
		log.Fatalf("Failed to write to .catppuccin.yaml")
	}
}

// MakeLocation saves the locations written to during installation into a file for later access.
func MakeLocation(packages string, location []string) {
	flavourrc := structs.AppLocation{
		Location: location,
	}
	marshallData, err := flavourrc.MarshalLocation()
	if err != nil {
		log.Error("Failed to marshall data.")
	}
	filepath := packages + ".yaml"
	finalPath := path.Join(ShareDir(), filepath)

	if PathExists(finalPath) { // If it already exists, remove it
		os.Remove(finalPath)
	}

	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Cannot open file.")
	}
	if _, err := file.Write(marshallData); err != nil {
		log.Error("Failed to write to file.")
	}
}
