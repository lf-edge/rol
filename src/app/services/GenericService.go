package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/dtos"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//GenericService generic service structure for IEntityModel
type GenericService[DtoType interfaces.IEntityDtoModel,
	CreateDtoType interfaces.IEntityDtoModel,
	UpdateDtoType interfaces.IEntityDtoModel,
	EntityType interfaces.IEntityModel] struct {
	repository    interfaces.IGenericRepository[EntityType]
	logger        *logrus.Logger
	logSourceName string
}

//NewGenericService generic service constructor
//Params
//	repo - generic repository
//	logger - logger
//Return
//	New generic service
func NewGenericService[DtoType interfaces.IEntityDtoModel,
	CreateDtoType interfaces.IEntityDtoModel,
	UpdateDtoType interfaces.IEntityDtoModel,
	EntityType interfaces.IEntityModel](repo interfaces.IGenericRepository[EntityType], logger *logrus.Logger) (
	*GenericService[
		DtoType,
		CreateDtoType,
		UpdateDtoType,
		EntityType], error) {
	model := new(EntityType)
	return &GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]{
		repository:    repo,
		logger:        logger,
		logSourceName: fmt.Sprintf("GenericService<%s>", reflect.TypeOf(*model).Name()),
	}, nil
}

func (s *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) log(ctx context.Context, level, message string) {
	entry := s.logger.WithFields(logrus.Fields{
		"actionId": ctx.Value("requestId"),
		"source":   s.logSourceName,
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

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) excludeDeleted(queryBuilder interfaces.IQueryBuilder) {
	date, err := time.Parse("2006-01-02", "1999-01-01")
	if err != nil {
		ges.logger.Errorf("time parse error: %s", err.Error())
		return
	}
	queryBuilder.Where("DeletedAt", "<", date)
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) addSearchInAllFields(search string, queryBuilder interfaces.IQueryBuilder) {
	entityModel := new(EntityType)
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(entityModel, stringFieldNames)
	queryGroup := ges.repository.NewQueryBuilder(nil)
	for i := 0; i < len(*stringFieldNames); i++ {
		fieldName := (*stringFieldNames)[i]
		containPass := strings.Contains(strings.ToLower(fieldName), "pass")
		containKey := strings.Contains(strings.ToLower(fieldName), "key")
		if containPass || containKey {
			continue
		}

		queryGroup.Or(fieldName, "LIKE", "%"+search+"%")
	}
	queryBuilder = queryBuilder.WhereQuery(queryGroup)
}

//GetList Get list of elements with filtering and pagination
//Params
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	*dtos.PaginatedListDto[DtoType] - pointer to paginated list
//	error - if an error occurs, otherwise nil
func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[DtoType], error) {
	pageFinal := page
	pageSizeFinal := pageSize
	if page < 1 {
		pageFinal = 1
	}
	if pageSize < 1 {
		pageSizeFinal = 10
	}
	searchQueryBuilder := ges.repository.NewQueryBuilder(ctx)
	ges.excludeDeleted(searchQueryBuilder)
	if len(search) > 3 {
		ges.addSearchInAllFields(search, searchQueryBuilder)
	}
	entities, err := ges.repository.GetList(ctx, orderBy, orderDirection, pageFinal, pageSizeFinal, searchQueryBuilder)
	if err != nil {
		return nil, err
	}
	count, err := ges.repository.Count(ctx, searchQueryBuilder)
	if err != nil {
		return nil, err
	}
	dtoArr := new([]DtoType)
	for i := 0; i < len(*entities); i++ {
		dto := new(DtoType)
		err = mappers.MapEntityToDto((*entities)[i], dto)
		if err != nil {
			return nil, fmt.Errorf("[%s]: [getList]: %s", ges.logSourceName, err.Error())
		}
		*dtoArr = append(*dtoArr, *dto)
	}

	paginatedDto := new(dtos.PaginatedListDto[DtoType])
	paginatedDto.Page = pageFinal
	paginatedDto.Size = pageSizeFinal
	paginatedDto.Total = count
	paginatedDto.Items = dtoArr
	return paginatedDto, nil
}

func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) getByIdExcludeDeleted(ctx context.Context, id uuid.UUID) (*EntityType, error) {
	queryBuilder := ges.repository.NewQueryBuilder(ctx)
	ges.excludeDeleted(queryBuilder)
	queryBuilder.Where("id", "=", id)
	entities, err := ges.repository.GetList(ctx, "id", "DESC", 1, 1, queryBuilder)
	if err != nil {
		return nil, err
	}
	if len(*entities) == 0 {
		return nil, nil
	}
	entity := &(*entities)[0]
	return entity, nil
}

//GetById Get entity by ID
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*DtoType - point to dto
//	error - if an error occurs, otherwise nil
func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetById(ctx context.Context, id uuid.UUID) (*DtoType, error) {
	entity, err := ges.getByIdExcludeDeleted(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	dto := new(DtoType)
	err = mappers.MapEntityToDto(*entity, dto)
	if err != nil {
		return nil, fmt.Errorf("[%s]: [getById]: %s", ges.logSourceName, err.Error())
	}
	return dto, nil
}

//Update save the changes to the existing entity
//Params
//	ctx - context is used only for logging
//	updateDto - update dto
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Update(ctx context.Context, updateDto UpdateDtoType, id uuid.UUID) error {
	entity, err := ges.getByIdExcludeDeleted(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		errStr := fmt.Sprintf("[%s]: [update]: entity with id %d does not exist", ges.logSourceName, id)
		err = errors.New(errStr)
		return err
	}
	err = mappers.MapDtoToEntity(updateDto, entity)
	if err != nil {
		return errors.New(fmt.Sprintf("[%s]: [update]: %s", ges.logSourceName, err.Error()))
	}
	err = ges.repository.Update(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

//Create add new entity
//Params
//	ctx - context is used only for logging
//	createDto - create dto
//Return
//	uuid.UUID - new entity id
//	error - if an error occurs, otherwise nil
func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Create(ctx context.Context, createDto CreateDtoType) (uuid.UUID, error) {
	entity := new(EntityType)
	err := mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return uuid.UUID{}, errors.New(fmt.Sprintf("[%s]: [update]: %s", ges.logSourceName, err.Error()))
	}
	id, err := ges.repository.Insert(ctx, *entity)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

//Delete mark entity as deleted
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (ges *GenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Delete(ctx context.Context, id uuid.UUID) error {
	entity, err := ges.getByIdExcludeDeleted(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		errStr := fmt.Sprintf("[%s]: [delete]: entity with id %d does not exist", ges.logSourceName, id)
		err = errors.New(errStr)
		return err
	}

	entityReflect := reflect.ValueOf(entity).Interface().(interfaces.IEntityModelDeletedAt)
	entityReflect.SetDeleted()

	err = ges.repository.Update(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
