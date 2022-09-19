package domain

//HostNetworkBridge is a struct describing network devices
type HostNetworkBridge struct {
	HostNetworkLink
	//Slaves slice of slaves interfaces names
	Slaves []string
}

//GetSlaves get bridge slaves
func (h HostNetworkBridge) GetSlaves() []string {
	return h.Slaves
}
