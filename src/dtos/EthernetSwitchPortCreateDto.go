package dtos

//EthernetSwitchPortCreateDto ethernet switch port create dto
type EthernetSwitchPortCreateDto struct {
	//	EthernetSwitchPortBaseDto - nested base switch port dto structure
	EthernetSwitchPortBaseDto
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esp EthernetSwitchPortCreateDto) Validate() error {
	return nil
}
