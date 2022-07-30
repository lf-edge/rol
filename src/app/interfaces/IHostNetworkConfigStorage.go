package interfaces

import "rol/domain"

type IHostNetworkConfigStorage interface {
	SaveConfig(config domain.HostNetworkConfig) error
	GetConfig() (domain.HostNetworkConfig, error)
	GetBackupConfig() (domain.HostNetworkConfig, error)
}
