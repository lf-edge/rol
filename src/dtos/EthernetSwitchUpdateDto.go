package dtos

type EthernetSwitchUpdateDto struct {
	EthernetSwitchBaseDto
	Password string
}

func (esd EthernetSwitchUpdateDto) Validate() error {
	return nil
}
