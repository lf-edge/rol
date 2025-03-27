package domain

import "github.com/google/uuid"

//DeviceNetworkInterface - device network interface configuration
type DeviceNetworkInterface struct {
	EntityUUID
	DeviceID uuid.UUID `gorm:"type:varchar(36);index"`
	//Mac address of network interface. This field is unique within all devices.
	Mac string `gorm:"type:varchar(17);index"`
	//ConnectedSwitchId - id of connected switch
	ConnectedSwitchId uuid.UUID `gorm:"type:varchar(36);index"`
	//ConnectedSwitchId - id of connected switch port
	ConnectedSwitchPortId uuid.UUID `gorm:"type:varchar(36);index"`
}
