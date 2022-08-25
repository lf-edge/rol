package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/mappers"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

func (e *EthernetSwitchService) portNameIsUniqueWithinTheSwitch(ctx context.Context, name string, switchID, id uuid.UUID) (bool, error) {
	uniqueNameQueryBuilder := e.portRepo.NewQueryBuilder(ctx)
	uniqueNameQueryBuilder.Where("Name", "==", name)
	uniqueNameQueryBuilder.Where("EthernetSwitchID", "==", switchID)
	if [16]byte{} != id {
		uniqueNameQueryBuilder.Where("ID", "!=", id)
	}
	ethSwitchPortsList, err := e.portRepo.GetList(ctx, "", "asc", 1, 1, uniqueNameQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "service failed to get list of switch ports")
	}
	if len(ethSwitchPortsList) > 0 {
		return false, nil
	}
	return true, nil
}

//GetPortByID Get ethernet switch port by switch ID and port ID
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetPortByID(ctx context.Context, switchID, id uuid.UUID) (dtos.EthernetSwitchPortDto, error) {
	dto := dtos.EthernetSwitchPortDto{}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return dto, errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	return GetByID[dtos.EthernetSwitchPortDto](ctx, e.portRepo, id, queryBuilder)
}

//CreatePort Create ethernet switch port by EthernetSwitchPortCreateDto
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	createDto - EthernetSwitchPortCreateDto
//Return
//	dtos.EthernetSwitchPortDto - created switch port
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) CreatePort(ctx context.Context, switchID uuid.UUID, createDto dtos.EthernetSwitchPortCreateDto) (dtos.EthernetSwitchPortDto, error) {
	dto := dtos.EthernetSwitchPortDto{}
	err := validators.ValidateEthernetSwitchPortCreateDto(createDto)
	if err != nil {
		return dto, err //we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return dto, errors.NotFound.New("ethernet switch is not found")
	}
	uniqName, err := e.portNameIsUniqueWithinTheSwitch(ctx, createDto.Name, switchID, [16]byte{})
	if err != nil {
		return dto, errors.Internal.Wrap(err, "name uniqueness check error")
	}
	if !uniqName {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Name", "port with this name already exist")
	}
	entity := new(domain.EthernetSwitchPort)
	entity.EthernetSwitchID = switchID
	err = mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	createdEntity, err := e.portRepo.Insert(ctx, *entity)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "create switch port in repository failed")
	}
	err = mappers.MapEntityToDto(createdEntity, &dto)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to map entity to dto")
	}
	return dto, nil
}

//UpdatePort Update ethernet switch port
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//  updateDto - dtos.EthernetSwitchPortUpdateDto DTO for updating entity
//Return
//	dtos.EthernetSwitchPortDto - updated switch port
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) UpdatePort(ctx context.Context, switchID, id uuid.UUID, updateDto dtos.EthernetSwitchPortUpdateDto) (dtos.EthernetSwitchPortDto, error) {
	dto := dtos.EthernetSwitchPortDto{}
	err := validators.ValidateEthernetSwitchPortUpdateDto(updateDto)
	if err != nil {
		return dto, err // we already wrap error in validators
	}
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return dto, errors.NotFound.New("switch is not found")
	}
	uniqName, err := e.portNameIsUniqueWithinTheSwitch(ctx, updateDto.Name, switchID, id)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "name uniqueness check error")
	}
	if !uniqName {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Name", "port with this name already exist")
	}
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	return Update[dtos.EthernetSwitchPortDto](ctx, e.portRepo, updateDto, id, queryBuilder)
}

//GetPorts Get list of ethernet switch ports with filtering and pagination
//Params
//	ctx - context is used only for logging
//  switchID - uuid of the ethernet switch
//	search - string for search in ethernet switch port string fields
//	orderBy - order by ethernet switch port field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedItemsDto[dtos.EthernetSwitchPortDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetPorts(ctx context.Context, switchID uuid.UUID, search, orderBy, orderDirection string,
	page, pageSize int) (dtos.PaginatedItemsDto[dtos.EthernetSwitchPortDto], error) {
	paginatedItemsDto := dtos.NewEmptyPaginatedItemsDto[dtos.EthernetSwitchPortDto]()
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return paginatedItemsDto, errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return paginatedItemsDto, errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	if len(search) > 3 {
		AddSearchInAllFields(search, e.portRepo, queryBuilder)
	}
	return GetListExtended[dtos.EthernetSwitchPortDto](ctx, e.portRepo, queryBuilder, orderBy, orderDirection, page, pageSize)
}

//DeletePort mark ethernet switch port as deleted
//Params
//	ctx - context is used only for logging
//	switchID - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) DeletePort(ctx context.Context, switchID, id uuid.UUID) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	_, err = e.portRepo.GetByIDExtended(ctx, id, queryBuilder)
	if err != nil {
		return err
	}
	err = e.portRepo.Delete(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to delete port")
	}
	return nil
}

func (e *EthernetSwitchService) deleteAllPortsBySwitchID(ctx context.Context, switchID uuid.UUID) error {
	switchExist, err := e.switchIsExist(ctx, switchID)
	if err != nil {
		return errors.Internal.Wrap(err, "error when checking the existence of the switch")
	}
	if !switchExist {
		return errors.NotFound.New("switch is not found")
	}
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	portsCount, err := e.portRepo.Count(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "ports counting failed")
	}
	ports, err := e.portRepo.GetList(ctx, "ID", "asc", 1, int(portsCount), queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get ports")
	}
	for _, port := range ports {
		err = e.portRepo.Delete(ctx, port.ID)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to remove port by id in repository")
		}
	}
	return nil
}
