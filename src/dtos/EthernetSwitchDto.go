package dtos

//EthernetSwitchDto
//	Ethernet switch response dto
type EthernetSwitchDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	BaseDto - nested base dto structure
	BaseDto
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (esd EthernetSwitchDto) Validate() error {
	return nil
}
