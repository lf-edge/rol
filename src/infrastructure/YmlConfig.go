package infrastructure

import (
	"os"
	"path"
	"rol/app/errors"
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
		return nil, errors.Internal.Wrapf(err, "failed to open config file %s", configFilePath)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, errors.Internal.Wrapf(err, "failed to parse yml config file %s", configFilePath)
	}
	return cfg, nil
}
