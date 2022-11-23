// Package domain stores the main structures of the program
package domain

//HostNetworkTrafficRule netfilter traffic rule
type HostNetworkTrafficRule struct {
	//Chain rule chain
	Chain string
	//Action rule action like ACCEPT, MASQUERADE, DROP, etc.
	Action string
	//Source packets source
	Source string
	//Destination packets destination
	Destination string
}
