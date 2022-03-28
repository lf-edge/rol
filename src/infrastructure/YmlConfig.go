package infrastructure

import (
	"fmt"
	"os"
	"path"
	"rol/domain"

	"gopkg.in/yaml.v2"
)

//NewYmlConfig create new config structure from yaml file
//Return
//	*domain.AppConfig - configuration structure
//	error - if error occurs return error, otherwise nil
func NewYmlConfig() (*domain.AppConfig, error) {
	cfg := &domain.AppConfig{}
	ex, _ := os.Executable()
	configFilePath := path.Join(path.Dir(ex), "appConfig.yml")
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("[NewYmlConfig]: Failed to open config file %s: %v", configFilePath, err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("[NewYmlConfig]: Failed to parse yml config file %s: %v", configFilePath, err)
	}
	return cfg, nil
}
