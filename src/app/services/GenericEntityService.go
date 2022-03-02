package services

import (
	"rol/app/interfaces"
	"rol/app/interfaces/generic"
	"rol/app/mappers"
)

type GenericEntityService struct {
	repository generic.IGenericEntityRepository
}

func NewGenericEntityService(repo generic.IGenericEntityRepository) (*GenericEntityService, error) {
	return &GenericEntityService{
		repository: repo,
	}, nil
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
