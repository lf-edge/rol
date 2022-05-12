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
)

type GenericRepositoryTest[EntityType interfaces.IEntityModel] struct {
	Repository interfaces.IGenericRepository[EntityType]
	Logger     *logrus.Logger
	Context    context.Context
	DbName     string
	InsertedId uuid.UUID
}

func NewGenericRepositoryTest[EntityType interfaces.IEntityModel](repo interfaces.IGenericRepository[EntityType], dbName string) *GenericRepositoryTest[EntityType] {
	return &GenericRepositoryTest[EntityType]{
		Repository: repo,
		Logger:     logrus.New(),
		Context:    context.TODO(),
		DbName:     dbName,
	}
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Insert(entity EntityType) error {
	var err error
	grt.InsertedId, err = grt.Repository.Insert(grt.Context, entity)
	if err != nil {
		return fmt.Errorf("insert failed: %s", err)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_GetById(id uuid.UUID) error {
	entity, err := grt.Repository.GetById(grt.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}

	value := reflect.ValueOf(*entity).FieldByName("ID")

	obtainedId, err := getUuidFromReflectArray(value)
	if obtainedId != id {
		return fmt.Errorf("unexpected id %d, expect %d", obtainedId, id)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Update(entity EntityType) error {
	err := grt.Repository.Update(grt.Context, &entity)

	if err != nil {
		return fmt.Errorf("update failed:  %s", err)
	}
	updatedEntity, err := grt.Repository.GetById(grt.Context, grt.InsertedId)
	if err != nil {
		return fmt.Errorf("get by id failed:  %s", err)
	}
	expectedName := reflect.ValueOf(entity).FieldByName("Name").String()
	obtainedName := reflect.ValueOf(*updatedEntity).FieldByName("Name").String()
	if obtainedName != expectedName {
		return fmt.Errorf("unexpected name %s, expect %s", obtainedName, expectedName)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_GetList() error {
	entityArr, err := grt.Repository.GetList(grt.Context, "", "", 1, 10, nil)
	if err != nil {
		return fmt.Errorf("get list failed:  %s", err)
	}
	if len(*entityArr) != 1 {
		return fmt.Errorf("array length %d, expect 1", len(*entityArr))
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Delete(id uuid.UUID) error {
	err := grt.Repository.Delete(grt.Context, id)
	if err != nil {
		return fmt.Errorf("delete failed:  %s", err)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Pagination(page, size int) error {
	entityArrFirstPage, err := grt.Repository.GetList(grt.Context, "created_at", "asc", page, size, nil)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArrFirstPage) != size {
		return fmt.Errorf("array length on %d page %d, expect %d", page, len(*entityArrFirstPage), size)
	}
	entityArrSecondPage, err := grt.Repository.GetList(grt.Context, "created_at", "asc", page+1, size, nil)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArrSecondPage) != size {
		return fmt.Errorf("array length on next page %d, expect %d", len(*entityArrSecondPage), size)
	}
	firstPageValue := reflect.ValueOf(*entityArrFirstPage).Index(0).FieldByName("ID")
	firstPageId, err := getUuidFromReflectArray(firstPageValue)
	if err != nil {
		return fmt.Errorf("convert reflect array to uuid failed: %s", err)
	}
	secondPageValue := reflect.ValueOf(*entityArrSecondPage).Index(0).FieldByName("ID")
	secondPageId, err := getUuidFromReflectArray(secondPageValue)
	if err != nil {
		return fmt.Errorf("convert reflect array to uuid failed: %s", err)
	}
	if firstPageId == secondPageId {
		return fmt.Errorf("pagination failed: got same element on second page with ID: %d", firstPageId)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Sort() error {
	entityArr, err := grt.Repository.GetList(grt.Context, "created_at", "desc", 1, 10, nil)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArr) != 10 {
		return fmt.Errorf("array length %d, expect 10", len(*entityArr))
	}
	index := len(*entityArr) / 2
	name := reflect.ValueOf(*entityArr).Index(index).FieldByName("Name").String()

	if name != fmt.Sprintf("AutoTesting_%d", index) {
		return fmt.Errorf("sort failed: got %s name, expect AutoTesting_%d", name, index)
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_Filter(queryBuilder interfaces.IQueryBuilder) error {
	entityArr, err := grt.Repository.GetList(grt.Context, "", "", 1, 10, queryBuilder)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArr) != 1 {
		return fmt.Errorf("array length %d, expect 1", len(*entityArr))
	}
	return nil
}

func (grt *GenericRepositoryTest[EntityType]) GenericRepository_CloseConnectionAndRemoveDb() error {
	if err := grt.Repository.CloseDb(); err != nil {
		return fmt.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove(grt.DbName); err != nil {
		return fmt.Errorf("remove db failed:  %s", err)
	}
	return nil
}

// getUuidFromReflectArray converts reflect.Value of array type to uuid.UUID
func getUuidFromReflectArray(value reflect.Value) (uuid.UUID, error) {
	bytes := make([]byte, 16)
	if value.Kind() == reflect.Array {
		// loop around each array element
		for i := 0; i < value.Len(); i++ {
			// get current element
			item := value.Index(i)
			// convert to uint type and record it to byte array
			u := item.Uint()
			bytes[i] = byte(u)
		}
		// create uuid from bytes array we created before
		id, err := uuid.FromBytes(bytes)
		if err != nil {
			return [16]byte{}, err
		}
		return id, nil
	} else {
		return [16]byte{}, errors.New("value is not a type of array")
	}
}
