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

//EthernetSwitchService service structure for EthernetSwitch entity
type EthernetSwitchService struct {
	*GenericService[dtos.EthernetSwitchDto,
		dtos.EthernetSwitchCreateDto,
		dtos.EthernetSwitchUpdateDto,
		domain.EthernetSwitch]
	supportedList *[]domain.EthernetSwitchModel
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
	if err != nil {
		return nil, err
	}
	ethernetSwitchService := &EthernetSwitchService{
		GenericService: genericService,
		supportedList:  &[]domain.EthernetSwitchModel{},
	}
	ethernetSwitchService.initSupportedList()
	return ethernetSwitchService, nil
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

func (e *EthernetSwitchService) sLog(ctx context.Context, level, message string) {
	entry := e.logger.WithFields(logrus.Fields{
		"actionID": ctx.Value("requestID"),
		"source":   "EthernetSwitchService",
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

func (e *EthernetSwitchService) modelIsSupported(model string) bool {
	modelIsSupported := false
	for _, supportedModel := range *e.supportedList {
		if model == supportedModel.Code {
			modelIsSupported = true
		}
	}
	return modelIsSupported
}

func (e *EthernetSwitchService) serialIsUnique(ctx context.Context, serial string, id uuid.UUID) error {
	uniqueSerialQueryBuilder := e.GenericService.repository.NewQueryBuilder(ctx)
	e.GenericService.excludeDeleted(uniqueSerialQueryBuilder)
	uniqueSerialQueryBuilder.Where("Serial", "==", serial)
	if [16]byte{} != id {
		uniqueSerialQueryBuilder.Where("ID", "!=", id)
	}
	serialEthSwitchList, err := e.GenericService.repository.GetList(ctx, "", "asc", 1, 1, uniqueSerialQueryBuilder)
	if err != nil {
		return fmt.Errorf("get list error: %s", err)
	}
	if len(*serialEthSwitchList) > 0 {
		return fmt.Errorf("switch with this serial number already exist")
	}
	return nil
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
	if !e.modelIsSupported(updateDto.SwitchModel) {
		return fmt.Errorf("switch model is not supported")
	}

	err = e.serialIsUnique(ctx, updateDto.Serial, id)
	if err != nil {
		return fmt.Errorf("serial number uniqueness check error: %s", err)
	}

	return e.GenericService.Update(ctx, updateDto, id)
}

//Create add new ethernet switch
//Params
//	ctx - context
//	createDto - ethernet switch create dto
//Return
//	uuid.UUID - new ethernet switch id
//	error - if an error occurs, otherwise nil
func (e *EthernetSwitchService) Create(ctx context.Context, createDto dtos.EthernetSwitchCreateDto) (uuid.UUID, error) {
	err := validators.ValidateEthernetSwitchCreateDto(createDto)
	if err != nil {
		return [16]byte{}, fmt.Errorf("dto validation error: %s", err)
	}
	if !e.modelIsSupported(createDto.SwitchModel) {
		return [16]byte{}, fmt.Errorf("this switch model is not supported")
	}
	err = e.serialIsUnique(ctx, createDto.Serial, [16]byte{})
	if err != nil {
		return [16]byte{}, fmt.Errorf("serial number uniqueness check error: %s", err)
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
