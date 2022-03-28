package dtos

//EthernetSwitchUpdateDto ethernet switch update dto
type EthernetSwitchUpdateDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	Password - ethernet switch management password
	Password string
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esd EthernetSwitchUpdateDto) Validate() error {
	return nil
}
