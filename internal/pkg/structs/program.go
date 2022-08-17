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
	InstallFiles    []string        `yaml:"install_files"`
	InstallFlavours InstallFlavours `yaml:"install_flavours"`
	OneFlavour      bool            `yaml:"one_flavour"`
	Modes           []string        `yaml:"modes"`
}

type InstallLocation struct {
	Unix string    `yaml:"unix"`
	Windows string `yaml:"windows"`
}

type Entry interface {}

type InstallFlavours struct {
	Latte     []Entry `yaml:"latte"`
	Frappe    []Entry `yaml:"frappe"`
	Macchiato []Entry `yaml:"macchiato"`
	Mocha     []Entry `yaml:"mocha"`
	To        string  `yaml:"to"`
}

