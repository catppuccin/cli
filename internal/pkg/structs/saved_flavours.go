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
	Location []string `yaml:"location"`
}
