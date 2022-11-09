package infrastructure

import (
	"os"
	"path"
	"rol/app/errors"
	"rol/domain"
)

//NewYmlConfig create new config structure from yaml file
//Return
//	*domain.AppConfig - configuration structure
//	error - if error occurs return error, otherwise nil
func NewYmlConfig() (*domain.AppConfig, error) {
	ex, _ := os.Executable()
	configFilePath := path.Join(path.Dir(ex), "appConfig.yml")
	cfg, err := ReadYamlFile[domain.AppConfig](configFilePath)
	if err != nil {
		return nil, errors.Internal.Wrapf(err, "failed to parse yml config file %s", configFilePath)
	}
	return &cfg, nil
}
