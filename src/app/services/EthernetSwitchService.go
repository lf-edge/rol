package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

//EthernetSwitchService service structure for EthernetSwitch entity
type EthernetSwitchService struct {
	switchRepo    interfaces.IGenericRepository[domain.EthernetSwitch]
	portRepo      interfaces.IGenericRepository[domain.EthernetSwitchPort]
	supportedList *[]domain.EthernetSwitchModel
}

//NewEthernetSwitchService constructor for domain.EthernetSwitch service
//Params
//	rep - generic repository with domain.EthernetSwitch repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchService(switchRepo interfaces.IGenericRepository[domain.EthernetSwitch],
	portRepo interfaces.IGenericRepository[domain.EthernetSwitchPort]) (interfaces.IGenericService[
	dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], error) {
	ethernetSwitchService := &EthernetSwitchService{
		switchRepo:    switchRepo,
		portRepo:      portRepo,
		supportedList: &[]domain.EthernetSwitchModel{},
	}
	return ethernetSwitchService, nil
}

//EthernetSwitchServiceInit do all that we need to do after dependency init
func EthernetSwitchServiceInit(service *EthernetSwitchService) error {
	service.initSupportedList()
	return nil
}

func (e *EthernetSwitchService) initSupportedList() {

	//Ubiquity UniFi Switch US-24-250W
	ubiquityUnifiSwitchUs24250W := domain.EthernetSwitchModel{
		Model:        "UniFi Switch US-24-250W",
		Manufacturer: "Ubiquity",
		Code:         "unifi_switch_us-24-250w",
	}
	*e.supportedList = append(*e.supportedList, ubiquityUnifiSwitchUs24250W)
}

func (e *EthernetSwitchService) modelIsSupported(model string) bool {
	modelIsSupported := false
	for _, supportedModel := range *e.supportedList {
		if model == supportedModel.Code {
			modelIsSupported = true
		}
	}
	return modelIsSupported
}

func (e *EthernetSwitchService) serialIsUnique(ctx context.Context, serial string, id uuid.UUID) (bool, error) {
	uniqueSerialQueryBuilder := e.switchRepo.NewQueryBuilder(ctx)
	uniqueSerialQueryBuilder.Where("Serial", "==", serial)
	if [16]byte{} != id {
		uniqueSerialQueryBuilder.Where("ID", "!=", id)
	}
	serialEthSwitchList, err := e.switchRepo.GetList(ctx, "", "asc", 1, 1, uniqueSerialQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "service failed get list")
	}
	if len(serialEthSwitchList) > 0 {
		return false, nil
	}
	return true, nil
}

func (e *EthernetSwitchService) addressIsUnique(ctx context.Context, serial string, id uuid.UUID) (bool, error) {
	uniqueSerialQueryBuilder := e.switchRepo.NewQueryBuilder(ctx)
	uniqueSerialQueryBuilder.Where("Address", "==", serial)
	if [16]byte{} != id {
		uniqueSerialQueryBuilder.Where("ID", "!=", id)
	}
	serialEthSwitchList, err := e.switchRepo.GetList(ctx, "", "asc", 1, 1, uniqueSerialQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "failed to get ethernet switches from repository")
	}
	if len(serialEthSwitchList) > 0 {
		return false, nil
	}
	return true, nil
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
//	dtos.PaginatedItemsDto[dtos.EthernetSwitchDto] - pointer to paginated list of ethernet switches
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.EthernetSwitchDto], error) {
	return GetList[dtos.EthernetSwitchDto](ctx, e.switchRepo, search, orderBy, orderDirection, page, pageSize)
}

//GetByID Get ethernet switch by ID
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	dtos.EthernetSwitchDto - point to ethernet switch dto
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) GetByID(ctx context.Context, id uuid.UUID) (dtos.EthernetSwitchDto, error) {
	return GetByID[dtos.EthernetSwitchDto](ctx, e.switchRepo, id, nil)
}

//Update save the changes to the existing ethernet switch
//Params
//	ctx - context is used only for logging
//	updateDto - ethernet switch update dto
//	id - ethernet switch id
//Return
//	dtos.EthernetSwitchDto - updated switch
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Update(ctx context.Context, updateDto dtos.EthernetSwitchUpdateDto, id uuid.UUID) (dtos.EthernetSwitchDto, error) {
	dto := dtos.EthernetSwitchDto{}
	err := validators.ValidateEthernetSwitchUpdateDto(updateDto)
	if err != nil {
		return dto, err // we already wrap error in validators
	}
	if !e.modelIsSupported(updateDto.SwitchModel) {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "SwitchModel", "this model is not supported")
	}

	uniqSerial, err := e.serialIsUnique(ctx, updateDto.Serial, id)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error occurred while checking uniqueness of the ethernet switch serial")
	}
	if !uniqSerial {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Serial", "ethernet switch with this serial number already exist")
	}

	uniqAddress, err := e.addressIsUnique(ctx, updateDto.Address, id)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "address uniqueness check error")
	}
	if !uniqAddress {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Address", "switch with this address already exist")
	}
	return Update[dtos.EthernetSwitchDto](ctx, e.switchRepo, updateDto, id, nil)
}

//Create add new ethernet switch
//Params
//	ctx - context
//	createDto - ethernet switch create dto
//Return
//	dtos.EthernetSwitchDto - created switch
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Create(ctx context.Context, createDto dtos.EthernetSwitchCreateDto) (dtos.EthernetSwitchDto, error) {
	dto := dtos.EthernetSwitchDto{}
	err := validators.ValidateEthernetSwitchCreateDto(createDto)
	if err != nil {
		return dtos.EthernetSwitchDto{}, errors.Validation.Wrap(err, "validation failed")
	}
	if !e.modelIsSupported(createDto.SwitchModel) {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "SwitchModel", "this model is not supported")
	}
	uniqSerial, err := e.serialIsUnique(ctx, createDto.Serial, [16]byte{})
	if err != nil {
		return dto, errors.Internal.Wrap(err, "error occurred while checking uniqueness of the ethernet switch serial")
	}
	if !uniqSerial {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Serial", "ethernet switch with this serial number already exist")
	}
	uniqAddress, err := e.addressIsUnique(ctx, createDto.Address, [16]byte{})
	if err != nil {
		return dto, errors.Validation.Wrap(err, "address uniqueness check error")
	}
	if !uniqAddress {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		return dto, errors.AddErrorContext(err, "Address", "switch with this address already exist")
	}
	dto, err = Create[dtos.EthernetSwitchDto](ctx, e.switchRepo, createDto)
	if err != nil {
		return dtos.EthernetSwitchDto{}, errors.Internal.Wrap(err, "service failed to create entity")
	}
	return dto, nil
}

//Delete mark ethernet switch as deleted
//Params
//	ctx - context is used only for logging
//	id - ethernet switch id
//Return
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Delete(ctx context.Context, id uuid.UUID) error {
	queryBuilder := e.portRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("EthernetSwitchID", "=", id)
	ports, err := e.portRepo.GetList(ctx, "", "", 1, 100, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get list of ports")
	}
	for _, port := range ports {
		err = e.portRepo.Delete(ctx, port.ID)
		if err != nil {
			return errors.Internal.Wrap(err, "service failed to update entity")
		}
	}
	err = e.switchRepo.Delete(ctx, id)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to delete entity from repository")
	}
	return nil
}

//GetSupportedModels Get supported switch models
//Return
//	*[]dtos.EthernetSwitchModelDto - Ethernet switch model DTO's that supported by system
func (e *EthernetSwitchService) GetSupportedModels() *[]dtos.EthernetSwitchModelDto {
	supportedModelsDtos := []dtos.EthernetSwitchModelDto{}
	for _, model := range *e.supportedList {
		modelDto := dtos.EthernetSwitchModelDto{}
		mappers.MapEthernetSwitchModelToDto(model, &modelDto)
		supportedModelsDtos = append(supportedModelsDtos, modelDto)
	}
	return &supportedModelsDtos
}
