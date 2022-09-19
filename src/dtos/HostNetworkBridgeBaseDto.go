package dtos

//HostNetworkBridgeBaseDto base dto for host network bridge
type HostNetworkBridgeBaseDto struct {
	//Addresses list
	Addresses []string
	//Slaves slice of slaves interfaces names
	Slaves []string
}
