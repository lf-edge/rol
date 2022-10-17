package domain

//DHCPServerState state for DHCP servers
type DHCPServerState uint

const (
	//DHCPStateLaunched dhcp v4 server in launched state
	DHCPStateLaunched = DHCPServerState(iota)
	//DHCPStateStopped dhcp v4 server in stopped state
	DHCPStateStopped
	//DHCPStateError dhcp v4 server in stopped state
	DHCPStateError
	//DHCPStateNone dhcp v4 server in not valid state
	DHCPStateNone
)

//String convert state to string
func (s DHCPServerState) String() string {
	switch s {
	case DHCPStateLaunched:
		return "launched"
	case DHCPStateStopped:
		return "stopped"
	case DHCPStateError:
		return "error"
	}
	return "unknown"
}
