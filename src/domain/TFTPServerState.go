// Package domain stores the main structures of the program
package domain

//TFTPServerState state for DHCP servers
type TFTPServerState uint

const (
	//TFTPStateLaunched tftp server in launched state
	TFTPStateLaunched = TFTPServerState(iota)
	//TFTPStateStopped tftp server in stopped state
	TFTPStateStopped
	//TFTPStateError tftp server in stopped state
	TFTPStateError
	//TFTPStateNone tftp server in unknown state
	TFTPStateNone
)

//String convert state to string
func (s TFTPServerState) String() string {
	switch s {
	case TFTPStateLaunched:
		return "launched"
	case TFTPStateStopped:
		return "stopped"
	case TFTPStateError:
		return "error"
	}
	return "unknown"
}
