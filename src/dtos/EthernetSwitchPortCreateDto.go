package dtos

type EthernetSwitchPortCreateDto struct {
	EthernetSwitchPortBaseDto
}

func (esp *EthernetSwitchPortCreateDto) Validate() error {
	return nil
}
