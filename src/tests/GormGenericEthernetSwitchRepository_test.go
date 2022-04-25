package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testRepoFileName string
var testRepoDbConnection gorm.Dialector
var testRepo interfaces.IGenericRepository[domain.EthernetSwitch]
var testCtx context.Context
var repoInsertedId uuid.UUID

func Test_GormGenericEthernetSwitchRepository_Prepare(t *testing.T) {
	testRepoFileName = "repository_test.db"
	testRepoDbConnection = sqlite.Open(testRepoFileName)
	testDb, err := gorm.Open(testRepoDbConnection, &gorm.Config{})
	err = testDb.AutoMigrate(
		&domain.EthernetSwitch{},
		&domain.EthernetSwitchPort{},
	)
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	testRepo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testDb, logger)
	testCtx = context.TODO()

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), testRepoFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(testRepoFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_GormGenericEthernetSwitchRepository_Insert(t *testing.T) {

	testEntity := domain.EthernetSwitch{
		Name:        "TestSwitch",
		Serial:      "TestSerial",
		Address:     "TestAddress",
		SwitchModel: 1,
		Username:    "TestUsername",
		Password:    "TestPassword",
	}

	var err error
	repoInsertedId, err = testRepo.Insert(testCtx, testEntity)
	if err != nil {
		t.Errorf("insert failed: %q", err)
	}
}

func Test_GormGenericEthernetSwitchRepository_GetById(t *testing.T) {
	testEntity, err := testRepo.GetById(testCtx, repoInsertedId)
	if err != nil {
		t.Errorf("get by id failed: %q", err)
	}
	if testEntity.Name != "TestSwitch" {
		t.Errorf("unexpected name %q, expect TestSwitch", testEntity.Name)
	}
}

func Test_GormGenericEthernetSwitchRepository_Update(t *testing.T) {
	newName := "UpdateTest"
	testEntity := &domain.EthernetSwitch{
		Name: newName,
	}
	testEntity.ID = repoInsertedId

	err := testRepo.Update(testCtx, testEntity)

	if err != nil {
		t.Errorf("update failed:  %q", err)
	}
	updatedEntity, err := testRepo.GetById(testCtx, repoInsertedId)
	if err != nil {
		t.Errorf("get by id failed:  %q", err)
	}
	if updatedEntity.Name != newName {
		t.Errorf("unexpected name %q, expect %q", updatedEntity.Name, newName)
	}
}

func Test_GormGenericEthernetSwitchRepository_GetList(t *testing.T) {
	entityArr, err := testRepo.GetList(testCtx, "", "", 1, 10, nil)
	if err != nil {
		t.Errorf("get list failed:  %q", err)
	}
	if len(*entityArr) != 1 {
		t.Errorf("array length %q, expect 1", len(*entityArr))
	}
}

func Test_GormGenericEthernetSwitchRepository_Delete(t *testing.T) {
	err := testRepo.Delete(testCtx, repoInsertedId)
	if err != nil {
		t.Errorf("delete failed:  %q", err)
	}
}

func Test_GormGenericEthernetSwitchRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		testEntity := domain.EthernetSwitch{
			Name:   fmt.Sprintf("TestSwitch_%d", i),
			Serial: fmt.Sprintf("%d", i),
		}
		_, err := testRepo.Insert(testCtx, testEntity)
		if err != nil {
			t.Errorf("insert failed:  %q", err)
		}
	}

	entityArr, err := testRepo.GetList(testCtx, "", "", 1, 20, nil)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*entityArr) != 20 {
		t.Errorf("array length %v, expect 20", len(*entityArr))
	}
}

func Test_GormGenericEthernetSwitchRepository_Pagination(t *testing.T) {
	entityArrFirstPage, err := testRepo.GetList(testCtx, "", "", 1, 10, nil)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*entityArrFirstPage) != 10 {
		t.Errorf("array length %v, expect 10", len(*entityArrFirstPage))
	}
	entityArrSecondPage, err := testRepo.GetList(testCtx, "", "", 2, 10, nil)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*entityArrSecondPage) != 10 {
		t.Errorf("array length %v, expect 10", len(*entityArrSecondPage))
	}
	searchId := (*entityArrSecondPage)[0].ID
	for _, entity := range *entityArrFirstPage {
		if entity.ID == searchId {
			t.Errorf("pagination failed: got same element on second page with ID: %q", searchId)
		}
	}
}

func Test_GormGenericEthernetSwitchRepository_Sort(t *testing.T) {
	entityArr, err := testRepo.GetList(testCtx, "created_at", "desc", 1, 10, nil)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*entityArr) != 10 {
		t.Errorf("array length %v, expect 10", len(*entityArr))
	}
	if (*entityArr)[5].Name != "TestSwitch_15" {
		t.Errorf("sort failed: got %s name, expect TestSwitch_15", (*entityArr)[5].Name)
	}
}

func Test_GormGenericEthernetSwitchRepository_Filter(t *testing.T) {
	queryBuilder := testRepo.NewQueryBuilder(testCtx)
	queryGroupBuilder := testRepo.NewQueryBuilder(testCtx)
	queryBuilder.Where("serial", "!=", "5").Where("serial", "!=", "9").WhereQuery(queryGroupBuilder.Where("name", "=", "TestSwitch_6").Or("name", "=", "TestSwitch_9"))
	entityArr, err := testRepo.GetList(testCtx, "", "", 1, 10, queryBuilder)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*entityArr) != 1 {
		t.Errorf("array length %v, expect 1", len(*entityArr))
	}
}

func Test_GormGenericEthernetSwitchRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	if err := testRepo.CloseDb(); err != nil {
		t.Errorf("close db failed:  %q", err)
	}
	if err := os.Remove(testRepoFileName); err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}
