package structs

import "gopkg.in/yaml.v3"

func UnmarshalLocation(data []byte) (AppLocation, error) {
	var r AppLocation
	err := yaml.Unmarshal(data, &r)
	return r, err
}

func (r *AppLocation) MarshalLocation() ([]byte, error) {
	return yaml.Marshal(r)
}

type AppLocation struct {
	AppName  string   `yaml:"app_name"`
	Location Location `yaml:"location"`
}

type Location struct {
	Directory string `yaml:"directory"`
}
