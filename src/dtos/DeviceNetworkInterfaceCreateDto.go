package dtos

import "github.com/google/uuid"

//DeviceNetworkInterfaceCreateDto - device network interface configuration create dto
type DeviceNetworkInterfaceCreateDto struct {
	//Mac address of network interface. This field is unique within all devices.
	Mac string
	//ConnectedSwitchId - id of connected switch
	ConnectedSwitchId uuid.UUID
	//ConnectedSwitchId - id of connected switch port
	ConnectedSwitchPortId uuid.UUID
}
