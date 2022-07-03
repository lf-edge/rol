package domain

//EthernetSwitch ethernet switch entity
type EthernetSwitch struct {
	//	Entity - nested base entity
	Entity
	//	Name - name of the switch
	Name string
	//	Serial - serial number of the switch
	Serial string
	//	SwitchModel - switch model code
	SwitchModel string
	//	Address - switch ip address
	Address string
	//	Username - switch management username
	Username string
	//	Password - switch management password
	Password string
	//	Ports - Switch ports
	Ports []EthernetSwitchPort
}
