package dtos

type EthernetSwitchCreateDto struct {
	EthernetSwitchBaseDto
	Password string
}

func (esd *EthernetSwitchCreateDto) Validate() error {
	return nil
}
