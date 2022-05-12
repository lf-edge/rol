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
)

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
	InsertedId uuid.UUID
}

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

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_Create(dto CreateDtoType) error {
	var err error
	gst.InsertedId, err = gst.Service.Create(gst.Context, dto)
	if err != nil {
		return fmt.Errorf("create failed: %s", err)
	}
	return nil
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_GetById(id uuid.UUID) error {
	dto, err := gst.Service.GetById(gst.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	if dto == nil {
		return fmt.Errorf("no entity with such id")
	}
	value := reflect.ValueOf(*dto).FieldByName("ID")

	obtainedId, err := getUuidFromReflectArray(value)
	if err != nil {
		return fmt.Errorf("convert bytes to uuid failed: %s", err)
	}
	if obtainedId != id {
		return fmt.Errorf("unexpected entity ID %d, expect %d", obtainedId, id)
	}
	return nil
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_Update(dto UpdateDtoType, id uuid.UUID) error {

	err := gst.Service.Update(gst.Context, dto, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	obtainedDto, err := gst.Service.GetById(gst.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}
	expectedName := reflect.ValueOf(dto).FieldByName("Name").String()
	obtainedName := reflect.ValueOf(*obtainedDto).FieldByName("Name").String()

	if obtainedName != expectedName {
		return fmt.Errorf("unexpected entity name %q, expect %q", obtainedName, expectedName)
	}
	return nil
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_Delete(id uuid.UUID) error {
	err := gst.Service.Delete(gst.Context, id)
	if err != nil {
		return fmt.Errorf("delete failed: %s", err)
	}
	dto, err := gst.Service.GetById(gst.Context, id)
	if dto != nil {
		return fmt.Errorf("unexpected entity %v, expect nil", dto)
	}
	return nil
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_GetList(total int64, page, size int) error {
	data, err := gst.Service.GetList(gst.Context, "", "CreatedAt", "desc", page, size)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if data.Total != total {
		return fmt.Errorf("get list failed: total items %d, expect %d", data.Total, total)
	}
	if data.Page != page {
		return fmt.Errorf("get list failed: page %d, expect %d", data.Page, page)
	}
	if data.Size != size {
		return fmt.Errorf("get list failed: size %d, expect %d", data.Size, size)
	}

	item := (*data.Items)[0]
	itemName := reflect.ValueOf(item).FieldByName("Name").String()
	expectedName := fmt.Sprintf("AutoTesting_%d", total)
	if itemName != expectedName {
		return fmt.Errorf("get list sort failed: get %s, expect AutoTesting_%d", itemName, total)
	}
	return nil
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_Search(search string) error {
	data, err := gst.Service.GetList(gst.Context, search, "", "", 1, 10)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*data.Items) < 1 {
		return errors.New("wasn't found any entries")
	}
	item := (*data.Items)[0]

	if containsInReflectStruct(item, search) {
		return nil
	}
	return fmt.Errorf("nothing was found for %s search word", search)
}

func (gst *GenericServiceTest[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GenericService_CloseConnectionAndRemoveDb() error {
	if err := gst.Repository.CloseDb(); err != nil {
		return fmt.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove(gst.DbName); err != nil {
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
