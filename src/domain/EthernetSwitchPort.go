package domain

import "github.com/google/uuid"

// EthernetSwitchPoePortType type of ethernet switch poe port
type EthernetSwitchPoePortType int

const (
	//POE standard poe
	POE = iota + 1
	//POE_PLUS poe plus
	POE_PLUS
	//PASSIVE_V24 passive v24 poe
	PASSIVE_V24
	//NONE none poe
	NONE = 10
)

//EthernetSwitchPort ethernet switch port entity
type EthernetSwitchPort struct {
	//	Entity - nested base entity
	Entity
	//	Name - name of switch port
	Name string
	//	EthernetSwitchID - id of switch port
	EthernetSwitchID uuid.UUID `gorm:"size:191"`
	//	PoeType - switch port type number
	PoeType EthernetSwitchPoePortType
}
