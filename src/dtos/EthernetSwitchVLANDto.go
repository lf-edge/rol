package dtos

import "github.com/google/uuid"

//EthernetSwitchVLANDto ethernet switch VLAN response dto
type EthernetSwitchVLANDto struct {
	BaseDto
	EthernetSwitchVLANBaseDto
	//VlanID VLAN ID
	VlanID int
	//EthernetSwitchID ethernet switch ID
	EthernetSwitchID uuid.UUID
}
