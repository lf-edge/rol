package dtos

import "rol/domain/enums"

type EthernetSwitchPortBaseDto struct {
	EthernetSwitchId uint
	PoeType          enums.EthernetSwitchPoePortType
	Name             string
}
