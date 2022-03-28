package dtos

//EthernetSwitchPortDto ethernet switch port response dto
type EthernetSwitchPortDto struct {
	//	EthernetSwitchPortBaseDto - nested base switch port dto structure
	EthernetSwitchPortBaseDto
	//	BaseDto - nested base dto structure
	BaseDto
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esp EthernetSwitchPortDto) Validate() error {
	return nil
}
