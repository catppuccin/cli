package structs

import (
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
}

type Catppuccinyaml struct {
	Name          string
	Exec          string
	MacosLocation string
	LinuxLocation string
	WinLocation   string
}
