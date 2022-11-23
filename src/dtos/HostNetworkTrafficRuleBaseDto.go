// Package dtos stores all data transfer objects
package dtos

//HostNetworkTrafficRuleBaseDto base dto for host network traffic rule
type HostNetworkTrafficRuleBaseDto struct {
	//Chain rule chain
	Chain string
	//Action rule action like ACCEPT, MASQUERADE, DROP, etc.
	Action string
	//Source packets source
	Source string
	//Destination packets destination
	Destination string
}
