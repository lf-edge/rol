package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/interfaces"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//EthernetSwitchService service structure for EthernetSwitch entity
type EthernetSwitchService struct {
	*GenericService[dtos.EthernetSwitchDto,
		dtos.EthernetSwitchCreateDto,
		dtos.EthernetSwitchUpdateDto,
		domain.EthernetSwitch]
}

//NewEthernetSwitchService constructor for domain.EthernetSwitch service
//Params
//	rep - generic repository with domain.EthernetSwitch repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchService(rep interfaces.IGenericRepository[domain.EthernetSwitch], log *logrus.Logger) (interfaces.IGenericService[
	dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], error) {
	genericService, err := NewGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch](rep, log)
	return &EthernetSwitchService{
		GenericService: genericService,
	}, err
}

//GetList Get list of ethernet switches with filtering and pagination
//Params
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedListDto[dtos.EthernetSwitchDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[dtos.EthernetSwitchDto], error) {
	return e.GenericService.GetList(ctx, search, orderBy, orderDirection, page, pageSize)
}

//GetByID Get ethernet switch by ID
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*dtos.EthernetSwitchDto - point to ethernet switch dto
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetByID(ctx context.Context, id uuid.UUID) (*dtos.EthernetSwitchDto, error) {
	return e.GenericService.GetByID(ctx, id)
}

//Update save the changes to the existing ethernet switch
//Params
//	ctx - context is used only for logging
//	updateDto - ethernet switch update dto
//	id - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Update(ctx context.Context, updateDto dtos.EthernetSwitchUpdateDto, id uuid.UUID) error {
	err := validators.ValidateEthernetSwitchUpdateDto(updateDto)
	if err != nil {
		return err
	}
	return e.GenericService.Update(ctx, updateDto, id)
}

//Create add new ethernet switch
//Params
//	ctx - context is used only for logging
//	createDto - ethernet switch create dto
//Return
//	uuid.UUID - new ethernet switch id
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Create(ctx context.Context, createDto dtos.EthernetSwitchCreateDto) (uuid.UUID, error) {
	err := validators.ValidateEthernetSwitchCreateDto(createDto)
	if err != nil {
		return [16]byte{}, err
	}
	return e.GenericService.Create(ctx, createDto)
}

//Delete mark ethernet switch as deleted
//Params
//	ctx - context is used only for logging
//	id - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Delete(ctx context.Context, id uuid.UUID) error {
	return e.GenericService.Delete(ctx, id)
}
