package structs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/browser"
	"gopkg.in/yaml.v3"
)

func UnmarshalProgram(data []byte) (Program, error) {
	var r Program
	err := yaml.Unmarshal(data, &r)
	return r, err
}

func (r *Program) Marshal() ([]byte, error) {
	return yaml.Marshal(r)
}

type Program struct {
	AppName      string       `yaml:"app_name"`
	BinaryName   string       `yaml:"binary_name"`
	Installation Installation `yaml:"installation"`
}

type InstallLocation struct {
	Linux   string `yaml:"linux"`
	Macos   string `yaml:"macos"`
	Windows string `yaml:"windows"`
}

type Entry struct {
	Default    []string            `yaml:"default"`
	Additional map[string][]string `yaml:"additional"`
}

type InstallFlavours struct {
	All       Entry `yaml:"all"`
	Latte     Entry `yaml:"latte"`
	Frappe    Entry `yaml:"frappe"`
	Macchiato Entry `yaml:"macchiato"`
	Mocha     Entry `yaml:"mocha"`
}

type Installation struct {
	InstallLocation InstallLocation `yaml:"location"`
	InstallFlavours InstallFlavours `yaml:"flavours"`
	To              string          `yaml:"to"`
	OneFlavour      bool            `yaml:"one_flavour"`
	Modes           []string        `yaml:"modes"`
	Comments        string          `yaml:"comments"`
	Hooks           Hooks           `yaml:"hooks"`
}


type Hooks struct {
  Pre  HookOptions `yaml:"post"`
  Post HookOptions `yaml:"post"`
}

type HookOptions struct {
	Install   []Hook `yaml:"install"`
	Uninstall []Hook `yaml:"uninstall"`
}

type HookType string

const (
	HookTypeShell   HookType = "shell"
	HookTypeBrowser HookType = "browser"
)

type Hook struct {
	Type HookType `yaml:"type"`
	Args []string `yaml:"args"`
}

func (h Hook) Run() error {
	if len(h.Args) == 0 {
		return errors.New("no args given for hook")
	}
	switch h.Type {
	case HookTypeShell:
		cmd := exec.Command(h.Args[0], h.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case HookTypeBrowser:
		return browser.OpenURL(h.Args[0])
	default:
		return fmt.Errorf("%q is an invalid hook type - .catppuccin.yaml invalid", h.Type)
	}
}

type Catppuccinyaml struct {
	Name          string
	Exec          string
	MacosLocation string
	LinuxLocation string
	WinLocation   string
}
