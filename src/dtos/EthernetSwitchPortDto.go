package dtos

import "github.com/google/uuid"

//EthernetSwitchPortDto ethernet switch port response dto
type EthernetSwitchPortDto struct {
	//	EthernetSwitchPortBaseDto - nested base switch port dto structure
	EthernetSwitchPortBaseDto
	//	BaseDto - nested base dto structure
	BaseDto[uuid.UUID]
}
