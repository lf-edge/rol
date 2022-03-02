package dtos

type EthernetSwitchDto struct {
	EthernetSwitchBaseDto
	BaseDto
}

func (esd *EthernetSwitchDto) Validate() error {
	return nil
}
