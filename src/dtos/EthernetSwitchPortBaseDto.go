package dtos

import "github.com/google/uuid"

//EthernetSwitchPortBaseDto base dto for ethernet switch port
type EthernetSwitchPortBaseDto struct {
	//	EthernetSwitchID - switch id
	EthernetSwitchID uuid.UUID
	// PoeType type of PoE for this port
	PoeType int
	// Name for this port
	Name string
}
