package dtos

//EthernetSwitchVLANCreateDto ethernet switch VLAN create dto
type EthernetSwitchVLANCreateDto struct {
	EthernetSwitchVLANBaseDto
	//VlanID VLAN ID
	VlanID int
}
