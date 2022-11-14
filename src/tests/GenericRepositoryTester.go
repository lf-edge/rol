package tests

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
	"testing"
	"time"
)

/////////////////////////////////////////////////////////////////
////////////////Entities and their fields section////////////////
/////////////////////////////////////////////////////////////////

//TestEntityFields fields for test entity
type TestEntityFields struct {
	String       string
	SecondString string
	SearchString string
	Yesterday    time.Time
	Tomorrow     time.Time
	NullDate     *time.Time
	NullableDate *time.Time
	Number       int
	Iterator     int
}

//Equals checks equivalence of th fields
func (f TestEntityFields) Equals(fields TestEntityFields) bool {
	if fields.NullDate != f.NullDate {
		return false
	}
	if fields.NullableDate != f.NullableDate {
		return false
	}
	if fields.String != f.String {
		return false
	}
	if fields.SearchString != f.SearchString {
		return false
	}
	if fields.SecondString != f.SecondString {
		return false
	}
	if fields.Number != f.Number {
		return false
	}
	if fields.Tomorrow != f.Tomorrow {
		return false
	}
	if fields.Yesterday != f.Yesterday {
		return false
	}
	return true
}

//GetTestEntityFields gets preset of entity test fields
func GetTestEntityFields() TestEntityFields {
	return GetTestEntityFieldsWithIteration(0)
}

//GetTestEntityFieldsWithIteration gets preset of entity test fields with selected iterator value
func GetTestEntityFieldsWithIteration(iteration int) TestEntityFields {
	nullableDate := time.Now()
	//DO NOT CHANGE!
	//IF YOU NEED ADD SOME NEW VALUES FOR NEW TEST ADD NEW FIELDS
	//AND DO NOT FORGET IMPLEMENT CHECK IN Equals() method
	return TestEntityFields{
		String:       "empty",
		SecondString: "second",
		SearchString: "search string with example text",
		NullDate:     nil,
		NullableDate: &nullableDate,
		Yesterday:    time.Now().AddDate(0, 0, -1),
		Tomorrow:     time.Now().AddDate(0, 0, 1),
		Number:       100,
		Iterator:     iteration,
	}
}

//BaseTestEntity base test entity struct (generic)
type BaseTestEntity[IDType comparable] struct {
	domain.Entity[IDType]
	TestEntityFields
}

//GetDataFields gets struct with entity data fields
func (t BaseTestEntity[IDType]) GetDataFields() TestEntityFields {
	return t.TestEntityFields
}

//TestEntityUUID base test entity struct with uuid id
type TestEntityUUID struct {
	BaseTestEntity[uuid.UUID]
}

//TestEntityInt base test entity struct with int id
type TestEntityInt struct {
	BaseTestEntity[int]
}

//ITestEntity default interface that we need for repository entity testing
type ITestEntity[IDType comparable] interface {
	interfaces.IEntityModel[IDType]
	GetDataFields() TestEntityFields
}

/////////////////////////////////////////////////////////////////
/////////////////// Generic Tester interface ////////////////////
/////////////////////////////////////////////////////////////////

//IGenericRepositoryTester interface that need to implement for each repository implementation
type IGenericRepositoryTester interface {
	SetTesting(testing *testing.T)
	Dispose() error
	TestInsertAndDelete()
	TestDeleteAll()
	TestGetByID()
	TestGetByIDExtended()
	TestDelete()
	TestIsExist()
	TestCount()
}

/////////////////////////////////////////////////////////////////
///////////////// Generic Tester implementation /////////////////
/////////////////////////////////////////////////////////////////

//GenericRepositoryTester generic struct for repository tester
type GenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]] struct {
	IDTypeName     string
	EntityTypeName string
	Implementation string
	ctx            context.Context
	t              *testing.T
	repo           interfaces.IGenericRepository[IDType, EntityType]
}

//NewGenericRepositoryTester constructor for generic repository tester
func NewGenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]](repo interfaces.IGenericRepository[IDType, EntityType]) (*GenericRepositoryTester[IDType, EntityType], error) {
	tester := &GenericRepositoryTester[IDType, EntityType]{}
	tester.EntityTypeName = reflect.TypeOf(*new(EntityType)).Name()
	tester.IDTypeName = reflect.TypeOf(*new(IDType)).Name()
	tester.repo = repo
	tester.ctx = context.Background()
	return tester, nil
}

func (t *GenericRepositoryTester[IDType, EntityType]) getTestBaseName() string {
	return fmt.Sprintf("%s/%s/%s", t.Implementation, t.EntityTypeName, t.IDTypeName)
}

//DataCleanUp clean up all data in repository
func (t *GenericRepositoryTester[IDType, EntityType]) dataCleanUp() {
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("CreatedAt", "!=", time.Now().AddDate(1, 1, 1))
	err := t.repo.DeleteAll(t.ctx, queryBuilder)
	if err != nil {
		panic("tests can't work correctly, database cleanup is not work")
	}
}

//SetTesting sets testing for interactive with tests stuff
func (t *GenericRepositoryTester[IDType, EntityType]) SetTesting(testing *testing.T) {
	t.t = testing
}

//Dispose all tester stuff
func (t *GenericRepositoryTester[IDType, EntityType]) Dispose() error {
	return t.repo.Dispose()
}

//createPredefinedTestEntityWithIteration create test entity with preset of fields for this repo with selected entity type
// with selected iterator value
func (t *GenericRepositoryTester[IDType, EntityType]) createPredefinedTestEntityWithIteration(iteration int) EntityType {
	entity := new(EntityType)
	var entityObj interface{} = entity
	if _, ok := entityObj.(interfaces.IEntityModel[uuid.UUID]); ok {
		entityObj = TestEntityUUID{BaseTestEntity: BaseTestEntity[uuid.UUID]{
			TestEntityFields: GetTestEntityFieldsWithIteration(iteration),
		}}
		return entityObj.(EntityType)
	}
	if _, ok := entityObj.(interfaces.IEntityModel[int]); ok {
		entityObj = TestEntityInt{BaseTestEntity: BaseTestEntity[int]{
			TestEntityFields: GetTestEntityFieldsWithIteration(iteration),
		}}
		return entityObj.(EntityType)
	}
	t.t.Error("Selected entity is not implement IEntityModel[IDType] interface")
	return *entity
}

//CreatePredefinedTestEntity create test entity with preset of fields for this repo with selected entity type
func (t *GenericRepositoryTester[IDType, EntityType]) createPredefinedTestEntity() EntityType {
	return t.createPredefinedTestEntityWithIteration(0)
}

//getNotExistedID gets not existed id with needed type
func (t *GenericRepositoryTester[IDType, EntityType]) getNotExistedID() IDType {
	id := new(IDType)
	var idObj interface{} = *id
	switch idObj.(type) {
	case int, int8, int16, int32, int64:
		idObj = 9999999
		return idObj.(IDType)
	case uuid.UUID:
		idObj = uuid.New()
		return idObj.(IDType)
	default:
		t.t.Error("getNotExistedID: id type is not supported")
	}
	return *id
}

func (t *GenericRepositoryTester[IDType, EntityType]) createEntitiesAndDeferClean(count int, f func(entities []EntityType)) func(test *testing.T) {
	return func(test *testing.T) {
		defer t.dataCleanUp()
		entities := []EntityType{}
		for i := 0; i < count; i++ {
			entity, err := t.repo.Insert(t.ctx, t.createPredefinedTestEntityWithIteration(i))
			if err != nil {
				t.t.Errorf("failed to create %d predefined entity", i)
			}
			entities = append(entities, entity)
		}
		if len(entities) != count {
			t.t.Errorf("failed to create %d predefined entities", count)
		}
		f(entities)
	}
}

//TestInsertAndDelete test inserting entity to repository and check that fields of entity wrote too
func (t *GenericRepositoryTester[IDType, EntityType]) TestInsertAndDelete() {
	t.t.Run(t.getTestBaseName(),
		t.createEntitiesAndDeferClean(0, func(entities []EntityType) {
			newEntity := t.createPredefinedTestEntity()
			newEntityFields := newEntity.GetDataFields()
			savedEntity, err := t.repo.Insert(t.ctx, newEntity)
			assert.NoError(t.t, err)
			savedEntityFields := savedEntity.GetDataFields()
			assert.Equal(t.t, true, savedEntityFields.Equals(newEntityFields))
			totalCount, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 1, totalCount)
			err = t.repo.Delete(t.ctx, savedEntity.GetID())
			assert.NoError(t.t, err)
		}))
}

//TestDeleteAll delete all test
func (t *GenericRepositoryTester[IDType, EntityType]) TestDeleteAll() {
	t.t.Run(t.getTestBaseName(),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("CreatedAt", "!=", time.Now().AddDate(1, 1, 1))
			err := t.repo.DeleteAll(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			totalCount, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 0, totalCount)
		}))
}

//TestGetByID get by id repo testing
func (t *GenericRepositoryTester[IDType, EntityType]) TestGetByID() {
	t.t.Run(fmt.Sprintf("%s/Found", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			//search
			getEntity, err := t.repo.GetByID(t.ctx, entities[1].GetID())
			assert.NoError(t.t, err)
			assert.Equal(t.t, getEntity.GetID(), entities[1].GetID())
		}))
	t.t.Run(fmt.Sprintf("%s/NotFound", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			_, err := t.repo.GetByID(t.ctx, t.getNotExistedID())
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
}

//TestGetByIDExtended get by id with conditions repo testing
func (t *GenericRepositoryTester[IDType, EntityType]) TestGetByIDExtended() {
	t.t.Run(fmt.Sprintf("%s/Found", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			getEntity, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, getEntity.GetID(), entities[1].GetID())
		}))
	t.t.Run(fmt.Sprintf("%s/NotFound", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			_, err := t.repo.GetByIDExtended(t.ctx, t.getNotExistedID(), nil)
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/EmtpyConditions", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
		}))
	t.t.Run(fmt.Sprintf("%s/NotExistedFieldCondition", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NotExistedField", "==", "empty")
			_, err := t.repo.GetByIDExtended(t.ctx, entities[2].GetID(), queryBuilder)
			assert.Error(t.t, err)
			assert.Equal(t.t, false, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/ExistedFieldCondition/Found", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "==", "second")
			receivedEntity, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, receivedEntity.GetID(), entities[1].GetID())
		}))
	t.t.Run(fmt.Sprintf("%s/FieldCondition/NotFound", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Number", "==", 1)
			_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/MultipleConditions/Found", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "==", "second").
				Where("SearchString", "LIKE", "%example%")
			_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.NotEqual(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/MultipleConditions/NotFound", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "!=", "second").
				Where("SearchString", "LIKE", "%example%")
			_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/ORMultipleConditions/Found", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "==", "second").
				Or("String", "!=", "empty")
			_, err := t.repo.GetByIDExtended(t.ctx, entities[0].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.NotEqual(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/ORMultipleConditions/NotFound", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Or("SecondString", "==", "notsecond").
				Or("String", "!=", "empty")
			_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
}

//TestDelete Deletion test
func (t *GenericRepositoryTester[IDType, EntityType]) TestDelete() {
	t.t.Run(fmt.Sprintf("%s/ExistedEntity", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			err := t.repo.Delete(t.ctx, entities[0].GetID())
			assert.NoError(t.t, err)
			_, err = t.repo.GetByID(t.ctx, entities[0].GetID())
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/NotExistedEntity", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			err := t.repo.Delete(t.ctx, t.getNotExistedID())
			assert.Error(t.t, err)
			assert.Equal(t.t, true, errors.As(err, errors.NotFound))
		}))
}

//TestIsExist testing for IsExist method
func (t *GenericRepositoryTester[IDType, EntityType]) TestIsExist() {
	t.t.Run(fmt.Sprintf("%s/NilConditions/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NilConditions/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			exist, err := t.repo.IsExist(t.ctx, t.getNotExistedID(), nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/EmptyConditions/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/EmptyConditions/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			exist, err := t.repo.IsExist(t.ctx, t.getNotExistedID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NotExistedFieldCondition", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NotExistedField", "==", "empty")
			_, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.Error(t.t, err)
			assert.Equal(t.t, false, errors.As(err, errors.NotFound))
		}))
	t.t.Run(fmt.Sprintf("%s/OneFieldCondition/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "==", "second")
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/OneFieldCondition/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Number", "==", 1)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/MultipleFieldsConditions/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "==", "second").
				Where("SearchString", "LIKE", "%example%")
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/MultipleFieldsConditions/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("SecondString", "!=", "second").
				Where("SearchString", "LIKE", "%example%")
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/ORMultipleFieldsConditions/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Or("SecondString", "==", "second").
				Or("String", "!=", "empty")
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/ORMultipleFieldsConditions/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Or("SecondString", "==", "notsecond").
				Or("String", "!=", "empty")
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NullFieldCondition/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NullDate", "==", nil)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NullFieldCondition/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NullDate", "!=", nil)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NullableFieldCondition/Exist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NullableDate", "!=", nil)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, true, exist)
		}))
	t.t.Run(fmt.Sprintf("%s/NullableFieldCondition/NotExist", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NullableDate", "==", nil)
			exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, false, exist)
		}))
}

//TestCount testing for Count method
func (t *GenericRepositoryTester[IDType, EntityType]) TestCount() {
	t.t.Run(fmt.Sprintf("%s/0Entity/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(0, func(entities []EntityType) {
			count, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 0, count)
		}))
	t.t.Run(fmt.Sprintf("%s/3Entity/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			count, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 3, count)
		}))
	t.t.Run(fmt.Sprintf("%s/20Entity/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(20, func(entities []EntityType) {
			count, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 20, count)
		}))
	t.t.Run(fmt.Sprintf("%s/NotExistedFieldCondition/ERROR", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("NotExistedField", "==", "empty")
			_, err := t.repo.Count(t.ctx, queryBuilder)
			assert.Error(t.t, err)
		}))
	t.t.Run(fmt.Sprintf("%s/ByIDOneElem/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Id", "==", entities[0].GetID())
			count, err := t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 1, count)
		}))
	t.t.Run(fmt.Sprintf("%s/ByIDTwoElem/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(3, func(entities []EntityType) {
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Id", "==", entities[0].GetID()).
				Or("Id", "==", entities[1].GetID())
			count, err := t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 2, count)
		}))
	t.t.Run(fmt.Sprintf("%s/IteratorRange/OK", t.getTestBaseName()),
		t.createEntitiesAndDeferClean(20, func(entities []EntityType) {
			count, err := t.repo.Count(t.ctx, nil)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 20, count)
			queryBuilder := t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Iterator", ">", 4).
				Where("Iterator", "<", 11)
			countIter, err := t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 6, countIter)
			queryBuilder = t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Iterator", ">=", 4).
				Where("Iterator", "<=", 11)
			countIterFirst, err := t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 8, countIterFirst)
			queryBuilder = t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Iterator", "<=", 4).
				Or("Iterator", ">=", 11)
			countIterSecond, err := t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 14, countIterSecond)
			queryBuilder = t.repo.NewQueryBuilder(t.ctx)
			queryBuilder.Where("Iterator", "<=", -20).
				Or("Iterator", ">=", 999)
			countIterSecond, err = t.repo.Count(t.ctx, queryBuilder)
			assert.NoError(t.t, err)
			assert.Equal(t.t, 0, countIterSecond)
		}))
}
