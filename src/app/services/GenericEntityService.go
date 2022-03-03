package services

import (
	"fmt"
	"regexp"
	"rol/app/interfaces"
	"rol/app/interfaces/generic"
	"rol/app/mappers"
	"rol/app/utils"
	"strings"
)

type GenericEntityService struct {
	repository generic.IGenericEntityRepository
}

func NewGenericEntityService(repo generic.IGenericEntityRepository) (*GenericEntityService, error) {
	return &GenericEntityService{
		repository: repo,
	}, nil
}

func toSnakeCase(entityStringFieldsNames *[]string) *[]string {
	snakeNames := &[]string{}
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	for i := 0; i < len(*entityStringFieldsNames); i++ {
		containPass := strings.Contains(strings.ToLower((*entityStringFieldsNames)[i]), "pass")
		containKey := strings.Contains(strings.ToLower((*entityStringFieldsNames)[i]), "key")
		if containPass || containKey {
			continue
		}
		snakeName := matchAllCap.ReplaceAllString((*entityStringFieldsNames)[i], "${1}_${2}")
		*snakeNames = append(*snakeNames, strings.ToLower(snakeName))
	}
	return snakeNames
}

func generateQuerySearchString(entity interface{}, search string) string {
	queryString := ""
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(entity, stringFieldNames)
	snakeCaseColumnNames := toSnakeCase(stringFieldNames)
	for i := 0; i < len(*snakeCaseColumnNames); i++ {
		if i != 0 {
			queryString = queryString + " OR "
		}
		queryString = queryString + fmt.Sprintf("%s LIKE \"%s\"", (*snakeCaseColumnNames)[i], "%"+search+"%")
	}
	return queryString
}

func (ges *GenericEntityService) GetList(dtoArr interface{}, search, orderBy, orderDirection string, page, pageSize int) (int64, error) {
	entities := mappers.GetEntityEmptyArray(dtoArr)
	queryRepString := ""
	var _ interface{} = nil
	if len(search) > 3 {
		entityModel := mappers.GetEmptyEntityFromArrayType(entities)
		queryRepString = generateQuerySearchString(entityModel, search)
	} else {
		queryRepString = ""
	}
	count, err := ges.repository.GetList(entities, orderBy, orderDirection, page, pageSize, queryRepString)
	mappers.Map(entities, dtoArr)
	return count, err
}

func (ges *GenericEntityService) GetAll(dtoArr interface{}) error {
	entities := mappers.GetEntityEmptyArray(dtoArr)
	ges.repository.GetAll(entities)
	mappers.Map(entities, dtoArr)
	return nil
}
func (ges *GenericEntityService) GetById(dto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(dto)
	err := ges.repository.GetById(entity, id)
	mappers.Map(entity, dto)
	if err != nil {
		return err
	}
	return nil
}
func (ges *GenericEntityService) Update(updateDto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(updateDto)
	err := ges.repository.GetById(entity, id)
	if err != nil {
		return err
	}
	mappers.Map(updateDto, entity)
	ges.repository.Update(entity)

	return nil
}
func (ges *GenericEntityService) Create(createDto interfaces.IEntityDtoModel) (uint, error) {
	entity := mappers.GetEmptyEntity(createDto)
	mappers.Map(createDto, entity)
	id, err := ges.repository.Insert(entity)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (ges *GenericEntityService) Delete(dto interfaces.IEntityDtoModel, id uint) error {
	entity := mappers.GetEmptyEntity(dto)
	err := ges.repository.GetById(entity, id)
	if err != nil {
		return err
	}
	return ges.repository.Delete(entity)
}
