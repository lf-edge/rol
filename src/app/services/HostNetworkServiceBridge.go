package services

import (
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
	"strings"
)

const slaveNotFound = "slave interface is not exist on the host"
const bridgeNotFound = "bridge is not exist on the host"

//GetBridgeList gets list of host bridges
//
//Return:
//	[]dtos.HostNetworkBridgeDto - slice of bridge dtos
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetBridgeList() ([]dtos.HostNetworkBridgeDto, error) {
	out := []dtos.HostNetworkBridgeDto{}
	links, err := h.manager.GetList()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting link list")
	}
	for _, link := range links {
		if link.GetType() == "bridge" && strings.Contains(link.GetName(), "rol.br.") {
			var dto dtos.HostNetworkBridgeDto
			err = mappers.MapEntityToDto(link, &dto)
			if err != nil {
				return nil, errors.Internal.Wrap(err, "error mapping bridge")
			}
			out = append(out, dto)
		}
	}
	return out, nil
}

//GetBridgeByName gets bridge by name
//
//Params:
//	name - name of the bridge
//Return:
//	dtos.HostNetworkBridgeDto - bridge dto
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetBridgeByName(name string) (dtos.HostNetworkBridgeDto, error) {
	out := dtos.HostNetworkBridgeDto{}
	link, err := h.manager.GetByName(name)
	if err != nil {
		return out, errors.Internal.Wrap(err, "error getting bridge by name")
	}
	if link == nil || link.GetType() != "bridge" || !strings.Contains(link.GetName(), "rol.br.") {
		return out, errors.NotFound.New(bridgeNotFound)
	}
	err = mappers.MapEntityToDto(link, &out)
	if err != nil {
		return out, errors.Internal.Wrap(err, "error mapping bridge")
	}
	return out, nil
}

func (h *HostNetworkService) syncBridgeSlaves(link interfaces.IHostNetworkLink, slaves []string) error {
	if link.GetType() != "bridge" {
		return errors.Internal.New("wrong link type received")
	}
	bridge := link.(domain.HostNetworkBridge)
	currSlaves := bridge.GetSlaves()
	deleteSlice, addSlice := utils.SliceDiffElements(currSlaves, slaves)
	for _, deleteSlave := range deleteSlice {
		err := h.manager.UnsetLinkMaster(deleteSlave)
		if err != nil {
			resetErr := h.manager.ResetChanges()
			if resetErr != nil {
				return errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, "set link no master failed")
		}
	}
	for _, addSlave := range addSlice {
		err := h.manager.SetLinkMaster(addSlave, bridge.GetName())
		if err != nil {
			resetErr := h.manager.ResetChanges()
			if resetErr != nil {
				return errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
			}
			return errors.Internal.Wrap(err, "set link master failed")
		}
	}
	return nil
}

//CreateBridge new bridge on host
//
//Params:
//	createDto - bridge create dto
//Return:
//	dtos.HostNetworkBridgeDto - created host network bridge
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) CreateBridge(createDto dtos.HostNetworkBridgeCreateDto) (dtos.HostNetworkBridgeDto, error) {
	dto := dtos.HostNetworkBridgeDto{}
	err := validators.ValidateHostNetworkBridgeCreateDto(createDto)
	if err != nil {
		return dto, err
	}
	for _, slave := range createDto.Slaves {
		_, err = h.manager.GetByName(slave)
		if err != nil {
			if !errors.As(err, errors.NotFound) {
				return dto, errors.Internal.Wrap(err, "failed to check existence of slave interface")
			}
			err1 := errors.Validation.New(errors.ValidationErrorMessage)
			return dto, errors.AddErrorContext(err1, "Slaves", slaveNotFound)
		}
	}
	bridgeName, err := h.manager.CreateBridge(createDto.Name)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error creating bridge")
	}
	err = h.manager.SetLinkUp(bridgeName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "set bridge up failed")
	}
	bridge, err := h.manager.GetByName(bridgeName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting bridge by name")
	}
	err = h.syncAddresses(bridge, createDto.Addresses)
	if err != nil {
		resetErr := h.manager.ResetChanges()
		if resetErr != nil {
			return dto, errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
		}
		return dto, err
	}
	err = h.syncBridgeSlaves(bridge, createDto.Slaves)
	if err != nil {
		resetErr := h.manager.ResetChanges()
		if resetErr != nil {
			return dto, errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
		}
		return dto, err
	}
	//Get updated bridge from manager
	bridge, err = h.manager.GetByName(bridgeName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting bridge by name")
	}
	err = mappers.MapEntityToDto(bridge, &dto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error mapping bridge")
	}
	return dto, nil
}

//UpdateBridge update bridge on host
//
//Params:
//	bridgeName - bridge name
//	updateDto - bridge update dto
//Return:
//	dtos.HostNetworkBridgeDto - updated bridge network vlan
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) UpdateBridge(bridgeName string, updateDto dtos.HostNetworkBridgeUpdateDto) (dtos.HostNetworkBridgeDto, error) {
	dto := dtos.HostNetworkBridgeDto{}
	err := validators.ValidateHostNetworkBridgeUpdateDto(updateDto)
	if err != nil {
		return dto, err
	}
	for _, slave := range updateDto.Slaves {
		_, err = h.manager.GetByName(slave)
		if err != nil {
			if !errors.As(err, errors.NotFound) {
				return dto, errors.Internal.Wrap(err, "failed to check existence of slave interface")
			}
			err1 := errors.Validation.New(errors.ValidationErrorMessage)
			return dto, errors.AddErrorContext(err1, "Slaves", slaveNotFound)
		}
	}
	bridge, err := h.manager.GetByName(bridgeName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting bridge by name")
	}
	if bridge == nil || (bridge.GetType() != "bridge" && strings.Contains(bridge.GetName(), "rol.br.")) {
		return dto, errors.NotFound.New("bridge not found")
	}
	err = h.syncAddresses(bridge, updateDto.Addresses)
	if err != nil {
		resetErr := h.manager.ResetChanges()
		if resetErr != nil {
			return dto, errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
		}
		return dto, err
	}
	err = h.syncBridgeSlaves(bridge, updateDto.Slaves)
	if err != nil {
		resetErr := h.manager.ResetChanges()
		if resetErr != nil {
			return dto, errors.Internal.Wrap(resetErr, "fatal: failed to reset changes after fail with setup address")
		}
		return dto, err
	}
	//Update link from manager
	bridge, err = h.manager.GetByName(bridgeName)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error getting bridge by name")
	}
	err = mappers.MapEntityToDto(bridge, &dto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error mapping bridge")
	}
	return dto, nil
}

//DeleteBridge deletes bridge on host by its name
//
//Params:
//	bridgeName - bridge name
//Return
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) DeleteBridge(bridgeName string) error {
	if !strings.Contains(bridgeName, "rol.br.") {
		return errors.NotFound.New(bridgeNotFound)
	}
	link, err := h.manager.GetByName(bridgeName)
	if err != nil {
		if !errors.As(err, errors.NotFound) {
			return errors.Internal.Wrap(err, "failed to check existence of bridgeName interface")
		}
		return err
	}
	if link == nil || link.GetType() != "bridge" {
		return errors.NotFound.New(bridgeNotFound)
	}
	err = h.manager.DeleteLinkByName(bridgeName)
	if err != nil {
		return errors.Internal.Wrap(err, "delete bridge failed")
	}
	return nil
}
