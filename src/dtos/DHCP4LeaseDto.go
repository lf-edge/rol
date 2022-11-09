package dtos

import (
	"github.com/google/uuid"
	"time"
)

//DHCP4LeaseDto DTO for DHCP v4 lease entity
type DHCP4LeaseDto struct {
	BaseDto[uuid.UUID]
	//IP address in ipv4 format
	IP string
	//MAC address in format like this 00-00-00-00-00
	MAC string
	//Expires datetime
	Expires time.Time
}
