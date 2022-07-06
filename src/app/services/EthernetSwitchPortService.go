package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//EthernetSwitchPortService service structure for EthernetSwitchPort entity
type EthernetSwitchPortService struct {
	*GenericService[dtos.EthernetSwitchPortDto,
		dtos.EthernetSwitchPortCreateDto,
		dtos.EthernetSwitchPortUpdateDto,
		domain.EthernetSwitchPort]
	switchRepository interfaces.IGenericRepository[domain.EthernetSwitch]
}

//NewEthernetSwitchPortService constructor for domain.EthernetSwitchPort service
//Params
//	rep - generic repository with domain.EthernetSwitchPort repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchPortService(rep interfaces.IGenericRepository[domain.EthernetSwitchPort], switchRepo interfaces.IGenericRepository[domain.EthernetSwitch], log *logrus.Logger) (interfaces.IGenericService[
	dtos.EthernetSwitchPortDto,
	dtos.EthernetSwitchPortCreateDto,
	dtos.EthernetSwitchPortUpdateDto,
	domain.EthernetSwitchPort], error) {
	genericService, err := NewGenericService[dtos.EthernetSwitchPortDto, dtos.EthernetSwitchPortCreateDto, dtos.EthernetSwitchPortUpdateDto, domain.EthernetSwitchPort](rep, log)
	if err != nil {
		return nil, err
	}
	ethernetSwitchService := &EthernetSwitchPortService{
		GenericService:   genericService,
		switchRepository: switchRepo,
	}
	return ethernetSwitchService, nil
}

func (e *EthernetSwitchPortService) switchExist(ctx context.Context, switchID uuid.UUID) error {
	queryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(queryBuilder)
	ethernetSwitch, err := e.switchRepository.GetByIDExtended(ctx, switchID, queryBuilder)
	if err != nil {
		return fmt.Errorf("failed to get ethernet switch: %s", err)
	}
	if ethernetSwitch == nil {
		return fmt.Errorf("failed to get ethernet switch: switch not found")
	}
	return nil
}

func (e *EthernetSwitchPortService) sLog(ctx context.Context, level, message string) {
	entry := e.logger.WithFields(logrus.Fields{
		"actionID": ctx.Value("requestID"),
		"source":   "EthernetSwitchPortService",
	})
	switch level {
	case "err":
		entry.Error(message)
	case "info":
		entry.Info(message)
	case "warn":
		entry.Warn(message)
	case "debug":
		entry.Debug(message)
	}
}

func (e *EthernetSwitchPortService) portNameIsUniqueWithinTheSwitch(ctx context.Context, name string, switchID, id uuid.UUID) error {
	uniqueNameQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueNameQueryBuilder)
	uniqueNameQueryBuilder.Where("Name", "==", name)
	uniqueNameQueryBuilder.Where("EthernetSwitchID", "==", switchID)
	if [16]byte{} != id {
		uniqueNameQueryBuilder.Where("ID", "!=", id)
	}
	ethSwitchPortsList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueNameQueryBuilder)
	if err != nil {
		return fmt.Errorf("get list error: %s", err)
	}
	if len(*ethSwitchPortsList) > 0 {
		return fmt.Errorf("switch with this address already exist")
	}
	return nil
}

//GetPortByID Get ethernet switch port by switch ID and port ID
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) GetPortByID(ctx context.Context, switchID, id uuid.UUID) (*dtos.EthernetSwitchPortDto, error) {
	err := e.switchExist(ctx, switchID)
	if err != nil {
		return nil, fmt.Errorf("error when checking the existence of the switch: %s", err)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchID", "==", switchID)
	return e.getByIDBasic(ctx, id, queryBuilder)
}

//CreatePort Create ethernet switch port by EthernetSwitchPortCreateDto
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	createDto - EthernetSwitchPortCreateDto
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) CreatePort(ctx context.Context, switchID uuid.UUID, createDto dtos.EthernetSwitchPortCreateDto) (uuid.UUID, error) {
	err := validators.ValidateEthernetSwitchPortCreateDto(createDto)
	if err != nil {
		return [16]byte{}, fmt.Errorf("dto validation error: %s", err)
	}
	err = e.switchExist(ctx, switchID)
	if err != nil {
		return [16]byte{}, fmt.Errorf("error when checking the existence of the switch: %s", err)
	}
	err = e.portNameIsUniqueWithinTheSwitch(ctx, createDto.Name, switchID, [16]byte{})
	if err != nil {
		return [16]byte{}, fmt.Errorf("name uniqueness check error: %s", err)
	}
	entity := new(domain.EthernetSwitchPort)
	entity.EthernetSwitchID = switchID
	err = mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to map ethernet switch port dto to entity: %s", err)
	}
	id, err := e.repository.Insert(ctx, *entity)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert ethernet switch port: %s", err)
	}
	return id, nil
}

//UpdatePort Update ethernet switch port
//Params
//	ctx - context is used only for logging|
//  switchID - Ethernet switch ID
//	id - entity id
//  updateDto - dtos.EthernetSwitchPortUpdateDto DTO for updating entity
//Return
//	*dtos.EthernetSwitchPortDto - point to ethernet switch port dto, if existed, otherwise nil
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) UpdatePort(ctx context.Context, switchID, id uuid.UUID, updateDto dtos.EthernetSwitchPortUpdateDto) error {
	err := validators.ValidateEthernetSwitchPortUpdateDto(updateDto)
	if err != nil {
		return fmt.Errorf("dto validation error: %s", err)
	}
	err = e.switchExist(ctx, switchID)
	if err != nil {
		return fmt.Errorf("error when checking the existence of the switch: %s", err)
	}
	err = e.portNameIsUniqueWithinTheSwitch(ctx, updateDto.Name, switchID, id)
	if err != nil {
		return fmt.Errorf("checking the uniqueness of the port name within one switch gave an error: %s", err)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	return e.updateBasic(ctx, updateDto, id, queryBuilder)
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
//	*dtos.PaginatedListDto[dtos.EthernetSwitchPortDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchPortService) GetPorts(ctx context.Context, switchID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[dtos.EthernetSwitchPortDto], error) {
	err := e.switchExist(ctx, switchID)
	if err != nil {
		return nil, fmt.Errorf("error when checking the existence of the switch: %s", err)
	}
	queryBuilder := e.repository.NewQueryBuilder(ctx)
	e.excludeDeleted(queryBuilder)
	queryBuilder.Where("EthernetSwitchId", "==", switchID)
	if len(search) > 3 {
		e.addSearchInAllFields(search, queryBuilder)
	}
	return e.getListBasic(ctx, queryBuilder, orderBy, orderDirection, page, pageSize)
}
