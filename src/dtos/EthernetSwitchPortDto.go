package dtos

type EthernetSwitchPortDto struct {
	EthernetSwitchPortBaseDto
	BaseDto
}

func (esp *EthernetSwitchPortDto) Validate() error {
	return nil
}
