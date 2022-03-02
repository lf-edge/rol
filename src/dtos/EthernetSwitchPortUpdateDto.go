package dtos

type EthernetSwitchPortUpdateDto struct {
	EthernetSwitchPortBaseDto
}

func (esp EthernetSwitchPortUpdateDto) Validate() error {
	return nil
}
