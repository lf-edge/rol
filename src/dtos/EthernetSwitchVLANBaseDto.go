package dtos

import "github.com/google/uuid"

//EthernetSwitchVLANBaseDto base dto for ethernet switch VLAN
type EthernetSwitchVLANBaseDto struct {
	//UntaggedPorts slice of untagged ports IDs
	UntaggedPorts []uuid.UUID
	//TaggedPorts slice of tagged ports IDs
	TaggedPorts []uuid.UUID
}
