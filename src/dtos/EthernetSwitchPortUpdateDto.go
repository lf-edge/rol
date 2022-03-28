package dtos

//EthernetSwitchPortUpdateDto ethernet switch port update dto
type EthernetSwitchPortUpdateDto struct {
	//	EthernetSwitchPortBaseDto - nested base switch port dto structure
	EthernetSwitchPortBaseDto
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esp EthernetSwitchPortUpdateDto) Validate() error {
	return nil
}
