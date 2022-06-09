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

//GenericRepositoryTest generic test for generic repository
type GenericRepositoryTest[EntityType interfaces.IEntityModel] struct {
	Repository interfaces.IGenericRepository[EntityType]
	Logger     *logrus.Logger
	Context    context.Context
	DbName     string
	InsertedID uuid.UUID
}

//NewGenericRepositoryTest GenericRepositoryTest constructor
func NewGenericRepositoryTest[EntityType interfaces.IEntityModel](repo interfaces.IGenericRepository[EntityType], dbName string) *GenericRepositoryTest[EntityType] {
	return &GenericRepositoryTest[EntityType]{
		Repository: repo,
		Logger:     logrus.New(),
		Context:    context.TODO(),
		DbName:     dbName,
	}
}

//GenericRepositoryInsert test insert entity to db
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryInsert(entity EntityType) error {
	var err error
	g.InsertedID, err = g.Repository.Insert(g.Context, entity)
	if err != nil {
		return fmt.Errorf("insert failed: %s", err)
	}
	return nil
}

//GenericRepositoryGetByID test get entity by ID
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryGetByID(id uuid.UUID) error {
	entity, err := g.Repository.GetByID(g.Context, id)
	if err != nil {
		return fmt.Errorf("get by id failed: %s", err)
	}

	value := reflect.ValueOf(*entity).FieldByName("ID")

	obtainedID, err := getUUIDFromReflectArray(value)
	if err != nil {
		return err
	}
	if obtainedID != id {
		return fmt.Errorf("unexpected id %d, expect %d", obtainedID, id)
	}
	return nil
}

//GenericRepositoryUpdate test update entity
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryUpdate(entity EntityType) error {
	err := g.Repository.Update(g.Context, &entity)

	if err != nil {
		return fmt.Errorf("update failed:  %s", err)
	}
	updatedEntity, err := g.Repository.GetByID(g.Context, g.InsertedID)
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

//GenericRepositoryGetList test get list of entities
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryGetList() error {
	entityArr, err := g.Repository.GetList(g.Context, "", "", 1, 10, nil)
	if err != nil {
		return fmt.Errorf("get list failed:  %s", err)
	}
	if len(*entityArr) != 1 {
		return fmt.Errorf("array length %d, expect 1", len(*entityArr))
	}
	return nil
}

//GenericRepositoryDelete test delete entity by ID
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryDelete(id uuid.UUID) error {
	err := g.Repository.Delete(g.Context, id)
	if err != nil {
		return fmt.Errorf("delete failed:  %s", err)
	}
	return nil
}

//GenericRepositoryPagination test pagination
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryPagination(page, size int) error {
	entityArrFirstPage, err := g.Repository.GetList(g.Context, "created_at", "asc", page, size, nil)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArrFirstPage) != size {
		return fmt.Errorf("array length on %d page %d, expect %d", page, len(*entityArrFirstPage), size)
	}
	entityArrSecondPage, err := g.Repository.GetList(g.Context, "created_at", "asc", page+1, size, nil)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArrSecondPage) != size {
		return fmt.Errorf("array length on next page %d, expect %d", len(*entityArrSecondPage), size)
	}
	firstPageValue := reflect.ValueOf(*entityArrFirstPage).Index(0).FieldByName("ID")
	firstPageID, err := getUUIDFromReflectArray(firstPageValue)
	if err != nil {
		return fmt.Errorf("convert reflect array to uuid failed: %s", err)
	}
	secondPageValue := reflect.ValueOf(*entityArrSecondPage).Index(0).FieldByName("ID")
	secondPageID, err := getUUIDFromReflectArray(secondPageValue)
	if err != nil {
		return fmt.Errorf("convert reflect array to uuid failed: %s", err)
	}
	if firstPageID == secondPageID {
		return fmt.Errorf("pagination failed: got same element on second page with ID: %d", firstPageID)
	}
	return nil
}

//GenericRepositorySort test sort
func (g *GenericRepositoryTest[EntityType]) GenericRepositorySort() error {
	entityArr, err := g.Repository.GetList(g.Context, "created_at", "desc", 1, 10, nil)
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

//GenericRepositoryFilter test filter
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryFilter(queryBuilder interfaces.IQueryBuilder) error {
	entityArr, err := g.Repository.GetList(g.Context, "", "", 1, 10, queryBuilder)
	if err != nil {
		return fmt.Errorf("get list failed: %s", err)
	}
	if len(*entityArr) != 1 {
		return fmt.Errorf("array length %d, expect 1", len(*entityArr))
	}
	return nil
}

//GenericRepositoryCloseConnectionAndRemoveDb close connection and remove db
func (g *GenericRepositoryTest[EntityType]) GenericRepositoryCloseConnectionAndRemoveDb() error {
	if err := g.Repository.CloseDb(); err != nil {
		return fmt.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove(g.DbName); err != nil {
		return fmt.Errorf("remove db failed:  %s", err)
	}
	return nil
}

//getUUIDFromReflectArray converts reflect.Value of array type to uuid.UUID
func getUUIDFromReflectArray(value reflect.Value) (uuid.UUID, error) {
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
	}
	return [16]byte{}, errors.New("value is not a type of array")
}
