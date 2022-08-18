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
	AppName         string          `yaml:"app_name"`
	PathName        string          `yaml:"path_name"`
	InstallLocation InstallLocation `yaml:"install_location"`
	InstallFlavours InstallFlavours `yaml:"install_flavours"`
	OneFlavour      bool            `yaml:"one_flavour"`
	Modes           []string        `yaml:"modes"`
}

type InstallLocation struct {
	Unix string    `yaml:"unix"`
	Windows string `yaml:"windows"`
}

type Additional interface {}

type Entry struct {
	Default    []string     `yaml:"default"`
	Additional Additional `yaml:"additional"` 
}

type InstallFlavours struct {
	All       Entry `yaml:"all"`
	Latte     Entry `yaml:"latte"`
	Frappe    Entry `yaml:"frappe"`
	Macchiato Entry `yaml:"macchiato"`
	Mocha     Entry `yaml:"mocha"`
	To        string  `yaml:"to"`
}

