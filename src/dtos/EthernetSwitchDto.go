package dtos

import "github.com/google/uuid"

//EthernetSwitchDto Ethernet switch response dto
type EthernetSwitchDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	BaseDto - nested base dto structure
	BaseDto[uuid.UUID]
}
