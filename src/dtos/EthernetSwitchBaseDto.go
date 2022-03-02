package dtos

import "rol/domain/enums"

type EthernetSwitchBaseDto struct {
	Name        string
	Serial      string
	SwitchModel enums.EthernetSwitchModel
	Address     string
	Username    string
}
