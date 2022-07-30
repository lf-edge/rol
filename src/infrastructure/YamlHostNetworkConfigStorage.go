package infrastructure

import (
	"os"
	"path/filepath"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//YamlHostNetworkConfigStorage Yaml implementation of IHostNetworkConfigStorage interface for host config management
type YamlHostNetworkConfigStorage struct {
	configFilePath string
}

//NewYamlHostNetworkConfigStorage constructor for YamlHostNetworkConfigStorage
func NewYamlHostNetworkConfigStorage(parameters domain.GlobalDIParameters) interfaces.IHostNetworkConfigStorage {
	return &YamlHostNetworkConfigStorage{
		configFilePath: filepath.Join(parameters.RootPath, "hostNetworkConfig.yaml"),
	}
}

//SaveConfig Save network configuration to config yaml file
func (y *YamlHostNetworkConfigStorage) SaveConfig(config domain.HostNetworkConfig) error {
	// Create backup config file
	if _, err := os.Stat(y.configFilePath); err == nil {
		backupFilePath := y.configFilePath + ".back"
		err = os.Rename(y.configFilePath, backupFilePath)
		if err != nil {
			return errors.Internal.Wrap(err, "error when creating backup host network config file")
		}
	}

	err := SaveYamlFile(config, y.configFilePath)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to save host network config file")
	}
	return nil
}

//GetConfig Gets configuration from config yaml file
func (y *YamlHostNetworkConfigStorage) GetConfig() (domain.HostNetworkConfig, error) {
	config, err := ReadYamlFile[domain.HostNetworkConfig](y.configFilePath)
	if err != nil {
		return domain.HostNetworkConfig{}, errors.NotFound.Wrap(err, "error reading host network config from file")
	}
	return config, nil
}

//GetBackupConfig Gets config from file config backup
func (y *YamlHostNetworkConfigStorage) GetBackupConfig() (domain.HostNetworkConfig, error) {
	config, err := ReadYamlFile[domain.HostNetworkConfig](y.configFilePath + ".back")
	if err != nil {
		return domain.HostNetworkConfig{}, errors.NotFound.Wrap(err, "backup of host network config is not found")
	}
	return config, nil
}
