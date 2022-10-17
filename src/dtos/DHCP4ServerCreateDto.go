package dtos

//DHCP4ServerCreateDto DTO for creating DHCP v4 server
type DHCP4ServerCreateDto struct {
	//Range of ip's for this dhcp v4 server, separated by "-", for example: "10.10.10.2-10.10.10.22"
	Range string
	//Mask for dhcp leases, for example: "255.255.255.0"
	Mask string
	//ServerID is a server_id dhcp option. Correct format is ipv4.
	ServerID string
	//Interface name
	Interface string
	//Gateway in ipv4 format
	Gateway string
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
