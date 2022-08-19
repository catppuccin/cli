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
	Properties   Properties   `yaml:"properties"`
	Installation Installation `yaml:"installation"`
	OneFlavour   bool         `yaml:"one_flavour"`
	Modes        []string     `yaml:"modes"`
}

type Properties struct {
	AppName         string          `yaml:"app_name"`
	BinaryName      string          `yaml:"binary_name"`
	InstallLocation InstallLocation `yaml:"install_location"`
}

type InstallLocation struct {
	Unix    string `yaml:"unix"`
	Windows string `yaml:"windows"`
}

type Additional interface{}

type Entry struct {
	Default    []string   `yaml:"default"`
	Additional Additional `yaml:"additional"`
}

type InstallFlavours struct {
	All       Entry `yaml:"all"`
	Latte     Entry `yaml:"latte"`
	Frappe    Entry `yaml:"frappe"`
	Macchiato Entry `yaml:"macchiato"`
	Mocha     Entry `yaml:"mocha"`
}

type Installation struct {
	InstallFlavours InstallFlavours `yaml:"flavours"`
	To              string          `yaml:"to"`
}
