package dtos

import "time"

//DHCP4LeaseCreateDto DTO for creating DHCP v4 lease entity
type DHCP4LeaseCreateDto struct {
	//IP address in ipv4 format
	IP string
	//MAC address in format like this 00-00-00-00-00
	MAC string
	//Expires datetime
	Expires time.Time
}
