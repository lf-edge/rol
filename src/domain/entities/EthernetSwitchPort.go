package entities

import (
	"rol/domain/base"
	"rol/domain/enums"
)

type EthernetSwitchPort struct {
	base.Entity
	Name             string
	EthernetSwitchId uint
	PoeType          enums.EthernetSwitchPoePortType
}
