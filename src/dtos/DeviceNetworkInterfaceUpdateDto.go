package dtos

import "github.com/google/uuid"

//DeviceNetworkInterfaceUpdateDto - device network interface configuration update dto
type DeviceNetworkInterfaceUpdateDto struct {
	//Mac address of network interface. This field is unique within all devices.
	Mac string
	//ConnectedSwitchId - id of connected switch
	ConnectedSwitchId uuid.UUID
	//ConnectedSwitchId - id of connected switch port
	ConnectedSwitchPortId uuid.UUID
}
