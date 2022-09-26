package dtos

//EthernetSwitchPortBaseDto base dto for ethernet switch port
type EthernetSwitchPortBaseDto struct {
	//POEType type of PoE for this port
	//can be: "poe", "poe+", "passive24", "none"
	POEType string
	// Name for this port
	Name string
	//POEEnabled POE status on a port
	POEEnabled bool
	//PVID default VLAN id
	PVID int
}
