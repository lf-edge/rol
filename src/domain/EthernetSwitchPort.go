package domain

import "github.com/google/uuid"

// EthernetSwitchPoePortType type of ethernet switch poe port
type EthernetSwitchPoePortType int

const (
	POE = iota + 1
	POE_PLUS
	PASSIVE_V24
	NONE = 10
)

//EthernetSwitchPort ethernet switch port entity
type EthernetSwitchPort struct {
	//	Entity - nested base entity
	Entity
	//	Name - name of switch port
	Name string
	//	EthernetSwitchId - id of switch port
	EthernetSwitchId uuid.UUID `gorm:"size:191"`
	//	PoeType - switch port type number
	PoeType EthernetSwitchPoePortType
}
