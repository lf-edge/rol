package dtos

//EthernetSwitchCreateDto ethernet switch create dto
type EthernetSwitchCreateDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	Password - ethernet switch management password
	Password string
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esd EthernetSwitchCreateDto) Validate() error {
	return nil
}
