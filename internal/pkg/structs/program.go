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
	InstallLocation string          `yaml:"install_location"`
	InstallFiles    []string        `yaml:"install_files"`
	InstallFlavours InstallFlavours `yaml:"install_flavours"`
	OneFlavour      bool            `yaml:"one_flavour"`
	Modes           []string        `yaml:"modes"`
}

type InstallFlavours struct {
	Latte     Flavour `yaml:"latte"`
	Frappe    Flavour `yaml:"frappe"`
	Macchiato Flavour `yaml:"macchiato"`
	Mocha     Flavour `yaml:"mocha"`
	To        string  `yaml:"to"`
}

type Entry interface {}

type Flavour struct {
	entries []Entry
}
