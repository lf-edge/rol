package services

import (
	"net"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/utils"
)

const parentNotFound = "parent interface is not exist on the host"
const setAddressesFailed = "set addresses to the interface fail"

//HostNetworkService is a struct for host network vlan service
type HostNetworkService struct {
	manager interfaces.IHostNetworkManager
}

//NewHostNetworkService is a constructor for HostNetworkService
//
//Params:
//	manager - host network manager
//Return:
//	HostNetworkService - instance of network vlan service
func NewHostNetworkService(manager interfaces.IHostNetworkManager) *HostNetworkService {
	return &HostNetworkService{
		manager: manager,
	}
}

func (h *HostNetworkService) syncAddresses(link interfaces.IHostNetworkLink, addresses []string) error {
	currAddresses := link.GetAddresses()
	linkName := link.GetName()
	currAddressStrSlice := []string{}
	for _, address := range currAddresses {
		currAddressStrSlice = append(currAddressStrSlice, address.String())
	}
	deletedCidrSlice, addedCidrSlice := utils.SliceDiffElements(currAddressStrSlice, addresses)
	for _, deletedCidr := range deletedCidrSlice {
		ip, address, err := net.ParseCIDR(deletedCidr)
		if err != nil {
			return errors.Internal.New("failed to parse CIDR")
		}
		address.IP = ip
		err = h.manager.AddrDelete(linkName, *address)
		if err != nil {
			err1 := h.manager.ResetChanges()
			if err1 != nil {
				return errors.Internal.Wrap(err, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, setAddressesFailed)
		}
	}
	for _, addedCidr := range addedCidrSlice {
		ip, address, err := net.ParseCIDR(addedCidr)
		if err != nil {
			return errors.Internal.New("failed to parse CIDR")
		}
		address.IP = ip
		err = h.manager.AddrAdd(linkName, *address)
		if err != nil {
			err1 := h.manager.ResetChanges()
			if err1 != nil {
				return errors.Internal.Wrap(err, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, setAddressesFailed)
		}
	}
	return nil
}
