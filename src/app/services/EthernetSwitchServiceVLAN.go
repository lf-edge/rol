package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

const (
	errorSwitchExistence = "error when checking the existence of the switch"
	errorPortExistence   = "error when checking the existence of the switch port"
	errorSwitchNotFound  = "switch is not found"
	errorGetPortByID     = "get port by id failed"
	errorAddTaggedVLAN   = "add tagged VLAN on port failed"
	errorRemoveVLAN      = "failed to remove VLAN from port"
	errorGetManager      = "can't get ethernet switch manager"
)

func (e *EthernetSwitchService) vlanIDIsUniqueWithinTheSwitch(ctx context.Context, vlanID int, switchID uuid.UUID) (bool, error) {
	uniqueIDQueryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	uniqueIDQueryBuilder.Where("VlanID", "==", vlanID)
	uniqueIDQueryBuilder.Where("EthernetSwitchID", "==", switchID)
	switchVLANsList, err := e.vlanRepo.GetList(ctx, "", "asc", 1, 1, uniqueIDQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "service failed to get list of switch ports")
	}
	if len(switchVLANsList) > 0 {
		return false, nil
	}
	return true, nil
}

func (e *EthernetSwitchService) checkNonexistentPorts(ctx context.Context, switchID uuid.UUID, dto dtos.EthernetSwitchVLANBaseDto) error {
	nonexistentTaggedPorts, err := e.getNonexistentPorts(ctx, switchID, dto.TaggedPorts)
	if err != nil {
		return errors.Internal.Wrap(err, errorPortExistence)
	}
	if len(nonexistentTaggedPorts) > 0 {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		for _, port := range nonexistentTaggedPorts {
			err = errors.AddErrorContext(err, "TaggedPorts", fmt.Sprintf("port %s doesn't exist", port.String()))
		}
		return err
	}

	nonexistentUntaggedPorts, err := e.getNonexistentPorts(ctx, switchID, dto.UntaggedPorts)
	if err != nil {
		return errors.Internal.Wrap(err, errorPortExistence)
	}
	if len(nonexistentUntaggedPorts) > 0 {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		for _, port := range nonexistentUntaggedPorts {
			err = errors.AddErrorContext(err, "UntaggedPorts", fmt.Sprintf("port %s doesn't exist", port.String()))
		}
		return err
	}
	return nil
}

func (e *EthernetSwitchService) checkExistenceOfRelatedEntities(ctx context.Context, switchID uuid.UUID, dto dtos.EthernetSwitchVLANBaseDto) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, errorSwitchExistence)
	}
	if !switchExist {
		return errors.NotFound.New(errorSwitchNotFound)
	}
	err = e.checkNonexistentPorts(ctx, switchID, dto)
	if err != nil {
		return err //we already wrap error
	}
	return nil
}

func (e *EthernetSwitchService) createVlanOnSwitch(ctx context.Context, switchID uuid.UUID, dto dtos.EthernetSwitchVLANCreateDto) error {
	switchManager, err := e.managers.Get(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, errorGetManager)
	}
	if switchManager == nil {
		return nil
	}
	err = switchManager.CreateVLAN(dto.VlanID)
	if err != nil {
		return errors.Internal.Wrap(err, "create VLAN on switch failed")
	}
	for _, taggedPortID := range dto.TaggedPorts {
		switchPort, err := e.portRepo.GetByID(ctx, taggedPortID)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.AddTaggedVLANOnPort(switchPort.Name, dto.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, errorAddTaggedVLAN)
		}
	}
	for _, untaggedPortID := range dto.UntaggedPorts {
		switchPort, err := e.portRepo.GetByID(ctx, untaggedPortID)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.AddUntaggedVLANOnPort(switchPort.Name, dto.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, errorAddTaggedVLAN)
		}
	}
	return nil
}

//GetVLANByID Get ethernet switch VLAN by switch ID and VLAN ID
//
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - VLAN id
//Return
//	*dtos.EthernetSwitchVLANDto - point to ethernet switch VLAN dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetVLANByID(ctx context.Context, switchID, id uuid.UUID) (dtos.EthernetSwitchVLANDto, error) {
	dto := dtos.EthernetSwitchVLANDto{}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, errorSwitchExistence)
	}
	if !switchExist {
		return dto, errors.NotFound.New(errorSwitchNotFound)
	}
	queryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	return GetByID[dtos.EthernetSwitchVLANDto, domain.EthernetSwitchVLAN](ctx, e.vlanRepo, id, queryBuilder)
}

//GetVLANs Get list of ethernet switch VLANs with filtering and pagination
//
//Params
//	ctx - context is used only for logging
//  switchID - uuid of the ethernet switch
//	search - string for search in ethernet switch port string fields
//	orderBy - order by ethernet switch port field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedListDto[dtos.EthernetSwitchVLANDto] - pointer to paginated list of ethernet switch VLANs
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetVLANs(ctx context.Context, switchID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.EthernetSwitchVLANDto], error) {
	dto := dtos.PaginatedItemsDto[dtos.EthernetSwitchVLANDto]{}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, errorSwitchExistence)
	}
	if !switchExist {
		return dto, errors.NotFound.New(errorSwitchNotFound)
	}
	queryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	if len(search) > 3 {
		AddSearchInAllFields(search, e.vlanRepo, queryBuilder)
	}
	return GetListExtended[dtos.EthernetSwitchVLANDto](ctx, e.vlanRepo, queryBuilder, orderBy, orderDirection, page, pageSize)
}

//CreateVLAN Create ethernet switch VLAN by EthernetSwitchVLANCreateDto
//Params
//	ctx - context is used only for logging
//  switchID - Ethernet switch ID
//	createDto - EthernetSwitchVLANCreateDto
//Return
//	uuid.UUID - created VLAN ID
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) CreateVLAN(ctx context.Context, switchID uuid.UUID, createDto dtos.EthernetSwitchVLANCreateDto) (dtos.EthernetSwitchVLANDto, error) {
	// Do all needed checks
	dto := dtos.EthernetSwitchVLANDto{}
	err := validators.ValidateEthernetSwitchVLANCreateDto(createDto)
	if err != nil {
		return dto, err //we already wrap error in validators
	}
	err = e.checkExistenceOfRelatedEntities(ctx, switchID, createDto.EthernetSwitchVLANBaseDto)
	if err != nil {
		return dto, err
	}
	uniqVLANId, err := e.vlanIDIsUniqueWithinTheSwitch(ctx, createDto.VlanID, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "VLAN ID uniqueness check error")
	}
	if !uniqVLANId {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "VlanID", "vlan with this id already exist")
	}

	// Create VLAN on Switch
	err = e.createVlanOnSwitch(ctx, switchID, createDto)
	if err != nil {
		return dto, err
	}

	// Save vlan configuration to repository
	entity := new(domain.EthernetSwitchVLAN)
	err = mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to map ethernet switch port dto to entity")
	}
	entity.EthernetSwitchID = switchID
	newVLAN, err := e.vlanRepo.Insert(ctx, *entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "repository failed to insert VLAN")
	}

	// Save configuration on the switch
	switchManager, err := e.managers.Get(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, errorGetManager)
	}
	if switchManager != nil {
		err = switchManager.SaveConfig()
		if err != nil {
			return dto, errors.Internal.Wrap(err, "save switch config failed")
		}
	}

	//Convert configuration to dto
	outDto := dtos.EthernetSwitchVLANDto{}
	err = mappers.MapEntityToDto(newVLAN, &outDto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to map vlan entity to dto")
	}
	return outDto, nil
}

//UpdateVLAN Update ethernet switch VLAN
//
//Params
//	ctx - context is used only for logging
//  switchID - Ethernet switch ID
//	id - VLAN id
//  updateDto - dtos.EthernetSwitchVLANUpdateDto DTO for updating entity
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) UpdateVLAN(ctx context.Context, switchID, id uuid.UUID, updateDto dtos.EthernetSwitchVLANUpdateDto) (dtos.EthernetSwitchVLANDto, error) {
	// Do all needed checks
	dto := dtos.EthernetSwitchVLANDto{}
	err := validators.ValidateEthernetSwitchVLANUpdateDto(updateDto)
	if err != nil {
		return dto, err //we already wrap error in validators
	}
	err = e.checkExistenceOfRelatedEntities(ctx, switchID, updateDto.EthernetSwitchVLANBaseDto)
	if err != nil {
		return dto, err
	}
	vlan, err := e.GetVLANByID(ctx, switchID, id)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "get VLAN by id failed")
	}

	// Apply all changes and save configuration on the switch
	err = e.syncVlanChangesOnSwitch(ctx, vlan, updateDto, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "update VLANs on port failed")
	}
	switchManager, err := e.managers.Get(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, errorGetManager)
	}
	if switchManager != nil {
		err = switchManager.SaveConfig()
		if err != nil {
			return dto, errors.Internal.Wrap(err, "save switch config failed")
		}
	}

	// Map entity to dto
	queryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	return Update[dtos.EthernetSwitchVLANDto](ctx, e.vlanRepo, updateDto, vlan.ID, queryBuilder)
}

//DeleteVLAN delete ethernet switch VLAN
//Params
//	ctx - context is used only for logging
//	switchID - ethernet switch id
//	id - VLAN ID
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) DeleteVLAN(ctx context.Context, switchID, id uuid.UUID) error {
	// Do all needed checks
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, errorSwitchExistence)
	}
	if !switchExist {
		return errors.NotFound.New(errorSwitchNotFound)
	}
	queryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	vlanEntity, err := e.vlanRepo.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return err
	}

	// Remove vlan and save configuration on switch
	switchManager, err := e.managers.Get(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, errorGetManager)
	}
	if switchManager != nil {
		err = switchManager.DeleteVLAN(vlanEntity.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to delete vlan from switch")
		}
		err = switchManager.SaveConfig()
		if err != nil {
			return errors.Internal.Wrap(err, "failed to save config")
		}
	}

	//Remove from repository
	return e.vlanRepo.Delete(ctx, id)
}

//syncVlanChangesOnSwitch apply vlan configuration to ethernet switch
func (e *EthernetSwitchService) syncVlanChangesOnSwitch(ctx context.Context, VLAN dtos.EthernetSwitchVLANDto, updateDto dtos.EthernetSwitchVLANUpdateDto, switchID uuid.UUID) error {
	switchManager, err := e.managers.Get(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, errorGetManager)
	}
	if switchManager == nil {
		return nil
	}
	diffTaggedToRemove, diffTaggedToAdd := utils.SliceDiffElements[uuid.UUID](VLAN.TaggedPorts, updateDto.TaggedPorts)
	for _, id := range diffTaggedToRemove {
		switchPort, err := e.portRepo.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.RemoveVLANFromPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, errorRemoveVLAN)
		}
	}
	for _, id := range diffTaggedToAdd {
		switchPort, err := e.portRepo.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.AddTaggedVLANOnPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to add tagged VLAN on port")
		}
	}
	diffUntaggedToRemove, diffUntaggedToAdd := utils.SliceDiffElements[uuid.UUID](VLAN.UntaggedPorts, updateDto.UntaggedPorts)
	for _, id := range diffUntaggedToRemove {
		switchPort, err := e.portRepo.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.RemoveVLANFromPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, errorRemoveVLAN)
		}
	}
	for _, id := range diffUntaggedToAdd {
		switchPort, err := e.portRepo.GetByID(ctx, id)
		if err != nil {
			return errors.Internal.Wrap(err, errorGetPortByID)
		}
		err = switchManager.AddUntaggedVLANOnPort(switchPort.Name, VLAN.VlanID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to add untagged VLAN on port")
		}
	}
	return nil
}

func (e *EthernetSwitchService) deleteAllVLANsBySwitchID(ctx context.Context, switchID uuid.UUID) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.vlanRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	vlansCount, err := e.vlanRepo.Count(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "VLANs counting failed")
	}
	vlans, err := e.vlanRepo.GetList(ctx, "ID", "asc", 1, int(vlansCount), queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get VLANs")
	}
	for _, vlan := range vlans {
		err = e.DeleteVLAN(ctx, switchID, vlan.ID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to remove one or more vlan that linked with this switch")
		}
	}
	return nil
}
