package structs

import (
  "gopkg.in/yaml.v3"
)

func UnmarshalFlavour(data []byte) (AppFlavour, error) {
  var r AppFlavour
  err := yaml.Unmarshal(data, &r)
  return r, err
}

func (r *AppFlavour) MarshalFlavour() ([]byte, error) {
  return yaml.Marshal(r)
}

type AppFlavour struct {
  AppName          string           `yaml:"app_name"`
  InstalledFlavour InstalledFlavour `yaml:"installed_flavour"`
}

type InstalledFlavour struct {
  Flavour string `yaml:"flavour"`
}
