package domain

import "github.com/google/uuid"

//EthernetSwitchVLAN ethernet switch VLAN entity
type EthernetSwitchVLAN struct {
	//EntityUUID - nested base entity where ID type is uuid.UUID
	EntityUUID
	VlanID           int
	EthernetSwitchID uuid.UUID `gorm:"index;size:36"`
	UntaggedPorts    string    `gorm:"type:text"`
	TaggedPorts      string    `gorm:"type:text"`
}
