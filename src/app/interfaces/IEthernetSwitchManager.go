package interfaces

//IEthernetSwitchManager is the interface is needed to manage ethernet switch
type IEthernetSwitchManager interface {
	//GetVLANs gets all VLANs on switch
	//
	//Return:
	//	[]int - slice of VLANs
	//	error - if an error occurs, otherwise nil
	GetVLANs() ([]int, error)
	//GetVLANsOnPort gets all VLANs on given port
	//
	//Params:
	//	portName - port name
	//Return:
	//	int - untagged VLAN ID
	//	[]int - slice of tagged VLANs IDs
	//	error - if an error occurs, otherwise nil
	GetVLANsOnPort(portName string) (int, []int, error)
	//AddTaggedVLANOnPort add tagged VLAN on given port
	//
	//Params:
	//	portName - port name
	//	vlanID - vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	AddTaggedVLANOnPort(portName string, vlanID int) error
	//AddUntaggedVLANOnPort add untagged VLAN on given port
	//
	//Params:
	//	portName - port name
	//	vlanID - vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	AddUntaggedVLANOnPort(portName string, vlanID int) error
	//RemoveVLANFromPort remove VLAN from given port
	//
	//Params:
	//	portName - name of port
	//	vlanID	- vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	RemoveVLANFromPort(portName string, vlanID int) error
	//SetPortPVID sets port PVID
	//
	//Params:
	//	portName - port name
	//	vlanID - vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	SetPortPVID(portName string, vlanID int) error
	//DeleteVLAN delete VLAN by id
	//
	//Params:
	//	vlanID - vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	DeleteVLAN(vlanID int) error
	//CreateVLAN create vlan on switch
	//
	//Params:
	//	vlanID	- vlan ID
	//Return:
	//	error - if an error occurs, otherwise nil
	CreateVLAN(vlanID int) error
	//GetPOEPortStatus gets poe status on given port
	//
	//Params:
	//	portName - port name
	//Return:
	//	string - poe port status "enable" or "disable"
	//	error - if an error occurs, otherwise nil
	GetPOEPortStatus(portName string) (string, error)
	//EnablePOEPort enable poe on give port
	//
	//Params:
	//	portName - port name
	//	poeType - poe type: "poe", "poe+" etc
	//Return:
	//	error - if an error occurs, otherwise nil
	EnablePOEPort(portName, poeType string) error
	//DisablePOEPort disable poe on given port
	//
	//Params:
	//	portName - port name
	//Return:
	//	error - if an error occurs, otherwise nil
	DisablePOEPort(portName string) error
	//SaveConfig save current settings on switch
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	SaveConfig() error
}
