package entities

import (
	"rol/domain/base"
	"rol/domain/enums"
)

type EthernetSwitch struct {
	base.Entity
	Name        string
	Serial      string
	SwitchModel enums.EthernetSwitchModel
	Address     string
	Username    string
	Password    string
	Ports       *[]*EthernetSwitchPort
}
