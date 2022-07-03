package dtos

import "github.com/google/uuid"

//EthernetSwitchPortBaseDto base dto for ethernet switch port
type EthernetSwitchPortBaseDto struct {
	//EthernetSwitchID - switch id
	EthernetSwitchID uuid.UUID
	//POEType type of PoE for this port
	//can be: "poe", "poe+", "passive24", "none"
	POEType string
	// Name for this port
	Name string
}
