package tests

import (
	"fmt"
	"reflect"
	"rol/app/interfaces"
	"rol/app/mappers"
)

//GenericMapperTest generic struct for test mapper
type GenericMapperTest[Dto interface{}, Entity interfaces.IEntityModel] struct{}

//NewGenericMapperToEntity GenericMapperTest constructor
func NewGenericMapperToEntity[Dto interface{}, Entity interfaces.IEntityModel]() *GenericMapperTest[Dto, Entity] {
	return &GenericMapperTest[Dto, Entity]{}
}

//MapToEntity maps dto to entity
func (g *GenericMapperTest[Dto, Entity]) MapToEntity(dto Dto, entity *Entity) error {
	err := mappers.MapDtoToEntity(dto, entity)
	if err != nil {
		return err
	}
	return nil
}

//MapToDto maps entity to dto
func (g *GenericMapperTest[Dto, Entity]) MapToDto(entity Entity, dto *Dto) error {
	err := mappers.MapEntityToDto(entity, dto)
	if err != nil {
		return err
	}
	return nil
}

func compareDtoAndEntity(from interface{}, to interface{}) error {
	fromReflectFields := reflect.TypeOf(from)
	fromReflectValues := reflect.ValueOf(from)

	num := fromReflectFields.NumField()
	for i := 0; i < num; i++ {
		reflectField := fromReflectFields.Field(i)
		reflectValue := fromReflectValues.Field(i)

		if reflectValue.Kind() == reflect.Struct && reflectField.Tag != "gorm:\"index\"" {
			if err := compareDtoAndEntity(reflectValue.Interface(), to); err != nil {
				return err
			}
		} else {
			if !reflectValue.IsValid() {
				continue
			}

			var fromFieldValue any
			if reflectValue.Kind() == reflect.Int {
				fromFieldValue = int(reflectValue.Int())
			} else {
				fromFieldValue = reflectValue.Interface()
			}

			toReflectValues := reflect.ValueOf(to)
			toValue := reflect.Indirect(toReflectValues).FieldByName(reflectField.Name)
			if !toValue.IsValid() {
				continue
			}

			var toFieldValue any
			if toValue.Kind() == reflect.Int {
				toFieldValue = int(toValue.Int())
			} else {
				toFieldValue = toValue.Interface()
			}

			if !reflect.DeepEqual(fromFieldValue, toFieldValue) {
				return fmt.Errorf("%s field not the same", reflectField.Name)
			}
		}
	}
	return nil
}
