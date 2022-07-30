package interfaces

import "rol/domain"

//IHostNetworkConfigStorage interface for network config management
type IHostNetworkConfigStorage interface {
	//SaveConfig Save host network configuration to storage
	//
	//Params
	//	config - configuration to save
	//Return
	//	error - if an error occurs, otherwise nil
	SaveConfig(config domain.HostNetworkConfig) error
	//GetConfig Get host network configuration from storage
	//
	//Return
	//	domain.HostNetworkConfig - configuration to save
	//	error - if an error occurs, otherwise nil
	GetConfig() (domain.HostNetworkConfig, error)
	//GetBackupConfig  Get backup of host network configuration from storage
	//
	//Return
	//	domain.HostNetworkConfig - configuration to save
	//	error - if an error occurs, otherwise nil
	GetBackupConfig() (domain.HostNetworkConfig, error)
}
