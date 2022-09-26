package domain

import "github.com/google/uuid"

//EthernetSwitchPort ethernet switch port entity
type EthernetSwitchPort struct {
	//Entity - nested base entity
	Entity
	//Name - name of switch port
	Name string
	//EthernetSwitchID - id of switch port
	EthernetSwitchID uuid.UUID `gorm:"size:191,index"`
	//POEType - switch port POE type
	//can be: "poe", "poe+", "passive24", "none"
	POEType string
	//POEEnabled POE status on a port
	POEEnabled bool
	//PVID - port PVID
	PVID int `gorm:"column:pvid"`
}
