package dtos

//EthernetSwitchUpdateDto ethernet switch update dto
type EthernetSwitchUpdateDto struct {
	//	EthernetSwitchBaseDto - nested base switch dto structure
	EthernetSwitchBaseDto
	//	Password - ethernet switch management password
	Password string
}
