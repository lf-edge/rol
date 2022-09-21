package infrastructure

import (
	"fmt"
	"regexp"
	"rol/app/errors"
	"rol/app/interfaces"
	"strconv"
	"strings"
)

const (
	//ErrorCreatingConnection error creating telnet connection
	ErrorCreatingConnection = "error creating telnet connection"
	//ErrorLoginIn error login in
	ErrorLoginIn = "error login in"
	//ErrorEnablingTelnet enabling failed
	ErrorEnablingTelnet = "enabling failed"
	//ErrorReadingTelnet error reading from telnet connection
	ErrorReadingTelnet = "error reading from telnet connection"
	//ErrorShowInterface show interface failed
	ErrorShowInterface = "show interface failed"
	//ErrorExecuteTelnet error executing telnet commands
	ErrorExecuteTelnet = "error executing telnet commands"
)

//TPLinkEthernetSwitchManager is a struct for tp link ethernet switch management
type TPLinkEthernetSwitchManager struct {
	telnetConn *TelnetConnection
	address    string
	login      string
	password   string
}

//NewTPLinkEthernetSwitchManager constructor for TPLinkEthernetSwitchManager
func NewTPLinkEthernetSwitchManager(address, login, password string) interfaces.IEthernetSwitchManager {
	tNet := NewTelnetConnection()
	return &TPLinkEthernetSwitchManager{
		telnetConn: tNet,
		address:    address,
		login:      login,
		//  pragma: allowlist nextline secret
		password: password,
	}
}

//GetVLANs gets all VLANs on switch
//
//Return:
//	[]int - slice of VLANs
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) GetVLANs() ([]int, error) {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return []int{}, errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return []int{}, errors.Internal.Wrap(err, ErrorLoginIn)
	}
	err = t.telnetConn.Send("enable")
	if err != nil {
		return []int{}, errors.Internal.Wrap(err, ErrorEnablingTelnet)
	}
	err = t.telnetConn.Send("show vlan")
	if err != nil {
		return []int{}, errors.Internal.Wrap(err, "showing vlan error")
	}
	_, err = t.telnetConn.Read("-----------\r\n")
	if err != nil {
		return []int{}, errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	out := []int{}
	for {
		msg, err := t.telnetConn.Read("\r")
		if err != nil {
			return []int{}, errors.Internal.Wrap(err, ErrorReadingTelnet)
		}
		fields := strings.Fields(msg)
		if len(fields) == 0 {
			break
		}
		if fields[0] == "System-VLAN" {
			id := 1
			out = append(out, id)
			continue
		}
		match, err := regexp.MatchString("[A-Z][a-z]\\d/\\d/\\d", fields[0])
		if err != nil {
			return nil, errors.Internal.Wrap(err, "regexp match string error")
		}
		if match {
			continue
		}
		if len(fields) > 1 {
			id, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, errors.Internal.Wrap(err, "error convert string to int")
			}
			out = append(out, id)
		}
	}
	return out, nil
}

//GetVLANsOnPort gets all ethernet switch VLANs on given port
//
//Params:
//	portName - port name
//Return:
//	int - untagged VLAN ID
//	[]int - slice of tagged VLANs IDs
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) GetVLANsOnPort(portName string) (int, []int, error) {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorLoginIn)
	}
	portNumber := portName[2:]
	err = t.telnetConn.Send("enable")
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorEnablingTelnet)
	}
	err = t.telnetConn.Send("show interface switchport gigabitEthernet " + portNumber)
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorShowInterface)
	}

	msg, err := t.telnetConn.Read("-----------\r\n")
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorReadingTelnet)
	}

	msg, err = t.telnetConn.Read("\r\n\n\r")
	if err != nil {
		return 0, []int{}, errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	untaggedVLAN := 0
	taggedVLANs := []int{}
	str := msg[:len(msg)-4]
	vlansInfo := strings.Split(str, "\r\n")
	for _, vlan := range vlansInfo {
		fields := strings.Fields(vlan)
		id, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, nil, errors.Internal.Wrap(err, "convert string to int failed")
		}
		if fields[2] == "Untagged" {
			untaggedVLAN = id
		} else {
			taggedVLANs = append(taggedVLANs, id)
		}
	}
	return untaggedVLAN, taggedVLANs, nil
}

//AddTaggedVLANOnPort add tagged VLAN on given port
//
//Params:
//	portName - port name
//	vlanID - vlan ID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) AddTaggedVLANOnPort(portName string, vlanID int) error {
	return t.addVLANOnPort(portName, "tagged", vlanID)
}

//RemoveVLANFromPort remove tagged VLAN from given port
//
//Params:
//	portName - name of port
//	vlanID	- vlan ID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) RemoveVLANFromPort(portName string, vlanID int) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	portNumber := portName[2:]
	exec := fmt.Sprintf("enable;config;interface gigabitEthernet %s;no switchport general allowed vlan %d;exit;exit;exit", portNumber, vlanID)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}

	return nil
}

//AddUntaggedVLANOnPort sets tagged VLAN on given port
//
//Params:
//	portName - port name
//	vlanID - vlan ID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) AddUntaggedVLANOnPort(portName string, vlanID int) error {
	return t.addVLANOnPort(portName, "untagged", vlanID)
}

//SetPortPVID set PVID on given port
//
//Params:
//	portName - port on which to set up the PVID
//	vlanID - PVID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) SetPortPVID(portName string, vlanID int) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	portNumber := portName[2:]
	exec := fmt.Sprintf("enable;config;interface gigabitEthernet %s;switchport pvid %d;end;exit;exit", portNumber, vlanID)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil
}

//DeleteVLAN delete VLAN by id
//
//Params:
//	vlanID - vlan ID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) DeleteVLAN(vlanID int) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	exec := fmt.Sprintf("enable;config;no vlan %d;exit;exit;exit", vlanID)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil
}

//CreateVLAN create vlan on switch
//
//Params:
//	vlanID	- vlan ID
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) CreateVLAN(vlanID int) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	exec := fmt.Sprintf("enable;config;vlan %d;exit;exit", vlanID)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil
}

//GetPOEPortStatus gets poe status on give port
//
//Params:
//	portName - port name
//Return:
//	string - poe port status "enable" or "disable"
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) GetPOEPortStatus(portName string) (string, error) {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorLoginIn)
	}
	portNumber := portName[2:]
	err = t.telnetConn.Send("enable")
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorEnablingTelnet)
	}
	err = t.telnetConn.Send("show power inline configuration interface gigabitEthernet " + portNumber)
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorShowInterface)
	}
	_, err = t.telnetConn.Read("-----------\r\n")
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	msg, err := t.telnetConn.Read("\r\n\n\r")
	if err != nil {
		return "", errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	portConfig := strings.Fields(msg)
	return portConfig[1], nil
}

//EnablePOEPort enable poe on given port
//
//Params:
//	portName - port name
//	poeType - poe type: "poe", "poe+" etc
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) EnablePOEPort(portName, poeType string) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	consumption := ""
	switch poeType {
	case "passive24":
		return errors.Internal.Wrap(err, "this switch does not support passive24 poe")
	default:
		consumption = "auto"
	}

	portNumber := portName[2:]
	exec := fmt.Sprintf("enable;config;interface gigabitEthernet %s;power inline consumption %s;power inline supply enable;exit;exit;exit;exit", portNumber, consumption)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil
}

//DisablePOEPort disable poe on given port
//
//Params:
//	portName - port name
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) DisablePOEPort(portName string) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	portNumber := portName[2:]
	exec := fmt.Sprintf("enable;config;interface gigabitEthernet %s;power inline supply disable;exit;exit;exit;exit", portNumber)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil
}

//SaveConfig Save current settings on switch
//
//Return:
//	error - if an error occurs, otherwise nil
func (t *TPLinkEthernetSwitchManager) SaveConfig() error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}
	exec := fmt.Sprintf("enable;copy running-config startup-config;exit;exit")
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	return nil

}

func (t *TPLinkEthernetSwitchManager) logIn() (err error) {
	_, err = t.telnetConn.Read("Login")
	if err != nil {
		return errors.Internal.Wrap(err, "error waiting for login string")
	}
	err = t.telnetConn.Send(t.login)
	if err != nil {
		return errors.Internal.Wrap(err, "login send error")
	}
	_, err = t.telnetConn.Read("Password")
	if err != nil {
		return errors.Internal.Wrap(err, "error waiting for password string")
	}
	err = t.telnetConn.Send(t.password)
	if err != nil {
		return errors.Internal.Wrap(err, "password send error")
	}
	return nil
}

func (t *TPLinkEthernetSwitchManager) addVLANOnPort(portName, vlanType string, vlanID int) error {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return errors.Internal.Wrap(err, ErrorLoginIn)
	}

	vlanExist, err := t.isVLANExists(vlanID)
	if err != nil {
		return errors.Internal.Wrap(err, "failed check vlan existence")
	}
	if vlanExist {
		portNumber := portName[2:]
		exec := fmt.Sprintf("enable;config;interface gigabitEthernet %s;switchport general allowed vlan %d %s;exit;exit;exit;", portNumber, vlanID, vlanType)
		err = t.executeTelnetCommands(exec)
		if err != nil {
			return errors.Internal.Wrap(err, ErrorExecuteTelnet)
		}
		return nil
	}
	return errors.NotFound.New("vlan not found")
}

func (t *TPLinkEthernetSwitchManager) isVLANExists(vlanID int) (bool, error) {
	err := t.telnetConn.Connect(t.address)
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorCreatingConnection)
	}
	err = t.logIn()
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorLoginIn)
	}
	err = t.telnetConn.Send("enable")
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorEnablingTelnet)
	}
	exec := fmt.Sprintf("config;show vlan id %d", vlanID)
	err = t.executeTelnetCommands(exec)
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorExecuteTelnet)
	}
	_, err = t.telnetConn.Read("---\r\n")
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	msg, err := t.telnetConn.Read("\r")
	if err != nil {
		return false, errors.Internal.Wrap(err, ErrorReadingTelnet)
	}
	vlan := strings.Fields(msg)
	if len(vlan) > 0 {
		return true, nil
	}
	return false, nil
}

func (t *TPLinkEthernetSwitchManager) executeTelnetCommands(exec string) error {
	commands := strings.Split(exec, ";")
	var err error
	for _, command := range commands {
		err = t.telnetConn.Send(command)
		if err != nil {
			return errors.Internal.Wrap(err, "send command to telnet server failed")
		}
	}
	return nil
}
