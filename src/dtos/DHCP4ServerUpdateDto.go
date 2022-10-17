package dtos

//DHCP4ServerUpdateDto DTO for updating DHCP v4 server
type DHCP4ServerUpdateDto struct {
	//DNS servers, separated by ";"
	DNS string
	//NTP IP address or dns name of NTP server
	NTP string
	//Enabled server or no
	Enabled bool
	//Port of DHCP server
	Port int
	//LeaseTime for dhcp v4 server leases
	LeaseTime int
}
