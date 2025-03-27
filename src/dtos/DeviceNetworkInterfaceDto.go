package dtos

import "github.com/google/uuid"

//DeviceNetworkInterfaceDto - device network interface configuration dto
type DeviceNetworkInterfaceDto struct {
	BaseDto[uuid.UUID]
	//Mac address of network interface. This field is unique within all devices.
	Mac string
	//ConnectedSwitchId - id of connected switch
	ConnectedSwitchId uuid.UUID
	//ConnectedSwitchId - id of connected switch port
	ConnectedSwitchPortId uuid.UUID
}
