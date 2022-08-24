package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"rol/app/interfaces"
	"strings"
	"time"
)

//GenericServiceTest generic test for generic service
type GenericServiceTest[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel] struct {
	Service interfaces.IGenericService[
		DtoType,
		CreateDtoType,
		UpdateDtoType,
		EntityType]
	Repository interfaces.IGenericRepository[EntityType]
	Logger     *logrus.Logger
	Context    context.Context
	DbName     string
	InsertedID uuid.UUID
}

//NewGenericServiceTest GenericServiceTest constructor
func NewGenericServiceTest[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel](
	service interfaces.IGenericService[DtoType,
		CreateDtoType,
		UpdateDtoType,
		EntityType],
	repo interfaces.IGenericRepository[EntityType], dbName string) *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType] {
	return &GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]{
		Service:    service,
		Repository: repo,
		Logger:     logrus.New(),
		Context:    context.TODO(),
		DbName:     dbName,
	}
}

//GenericServiceCreate test create entity
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceCreate(dto CreateDtoType) (DtoType, error) {
	return g.Service.Create(g.Context, dto)
}

//GenericServiceGetByID test get entity by id
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceGetByID(id uuid.UUID) error {
	dto, err := g.Service.GetByID(g.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	value := reflect.ValueOf(dto).FieldByName("ID")

	obtainedID, err := getUUIDFromReflectArray(value)
	if err != nil {
		return fmt.Errorf("convert bytes to uuid failed: %s", err)
	}
	if obtainedID != id {
		return fmt.Errorf("unexpected entity ID %d, expect %d", obtainedID, id)
	}
	return nil
}

//GenericServiceUpdate test update entity
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceUpdate(dto UpdateDtoType, id uuid.UUID) error {
	entity, err := g.Repository.GetByID(g.Context, g.InsertedID)
	if err != nil {
		return err
	}
	beforeUpdTime := reflect.ValueOf(entity).FieldByName("UpdatedAt").Interface().(time.Time)
	_, err = g.Service.Update(g.Context, dto, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	obtainedDto, err := g.Service.GetByID(g.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	afterUpdTime := reflect.ValueOf(obtainedDto).FieldByName("UpdatedAt").Interface().(time.Time)
	if !beforeUpdTime.Before(afterUpdTime) {
		return fmt.Errorf("entity was not updated")
	}
	return nil
}

//GenericServiceDelete test delete entity
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceDelete(id uuid.UUID) error {
	return g.Service.Delete(g.Context, id)
}

//GenericServiceGetList test get list of entities
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceGetList(total, page, size int) error {
	data, err := g.Service.GetList(g.Context, "", "CreatedAt", "desc", page, size)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if data.Pagination.TotalCount != total {
		return fmt.Errorf("get list failed: total items %d, expect %d", data.Pagination.TotalCount, total)
	}
	if data.Pagination.Page != page {
		return fmt.Errorf("get list failed: page %d, expect %d", data.Pagination.Page, page)
	}
	if data.Pagination.Size != size {
		return fmt.Errorf("get list failed: size %d, expect %d", data.Pagination.Size, size)
	}

	item := data.Items[0]
	itemName := reflect.ValueOf(item).FieldByName("Name").String()
	expectedName := fmt.Sprintf("AutoTesting_%d", total)
	if itemName != expectedName {
		return fmt.Errorf("get list sort failed: get %s, expect AutoTesting_%d", itemName, total)
	}
	return nil
}

//GenericServiceSearch test get list search
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceSearch(search string) error {
	data, err := g.Service.GetList(g.Context, search, "", "", 1, 10)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(data.Items) < 1 {
		return errors.New("wasn't found any entries")
	}
	item := data.Items[0]

	if containsInReflectStruct(item, search) {
		return nil
	}
	return fmt.Errorf("nothing was found for %s search word", search)
}

//GenericServiceCloseConnectionAndRemoveDb close connection and remove db
func (g *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericServiceCloseConnectionAndRemoveDb() error {
	if err := g.Repository.CloseDb(); err != nil {
		return fmt.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove(g.DbName); err != nil {
		return fmt.Errorf("remove db failed:  %s", err)
	}
	return nil
}

func containsInReflectStruct(item interface{}, search string) bool {
	// get fields reflect.Types and reflect.Values
	fields := reflect.TypeOf(item)
	values := reflect.ValueOf(item)
	num := fields.NumField()
	// loop around each field
	for i := 0; i < num; i++ {
		value := values.Field(i)
		// if field is a struct make recursion call for that field
		if value.Kind() == reflect.Struct {
			if containsInReflectStruct(value.Interface(), search) {
				return true
			}
			// if field is a string, look for a substring in it
		} else if value.Kind() == reflect.String {
			val := value.String()
			found := strings.Contains(val, search)
			if found {
				return true
			}
		}
	}
	return false
}
