package dtos

//EthernetSwitchCreateDto ethernet switch create dto
type EthernetSwitchCreateDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	Password - ethernet switch management password
	Password string
}
