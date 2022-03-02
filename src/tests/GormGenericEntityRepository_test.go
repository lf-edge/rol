package testss

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"os"
	"path"
	"rol/domain/base"
	"rol/domain/entities"
	"rol/infrastructure"
	"runtime"
	"testing"
)

var testRepoFileName = "repository_test.db"
var testRepoDbConnection = sqlite.Open(testRepoFileName)
var testRepo, _ = infrastructure.NewGormGenericEntityRepository(testRepoDbConnection)
var repoTestInsertedId uint = 0

func Test_GormGenericEntityRepository_Prepare(t *testing.T) {
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), testRepoFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err := os.Remove(testRepoFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_GormGenericEntityRepository_Insert(t *testing.T) {
	entity := &entities.EthernetSwitch{
		Entity:      base.Entity{},
		Name:        "AutoTesting",
		Serial:      "AutoTesting",
		SwitchModel: 0,
		Address:     "10.10.88.2",
		Username:    "AutoTesting",
		Password:    "AutoTesting",
	}
	repoTestInsertedId, _ = testRepo.Insert(entity)
	if repoTestInsertedId == 0 {
		t.Errorf("got %q, wanted %q", repoTestInsertedId, " > 0")
	}
	if entity.Name != "AutoTesting" {
		t.Errorf("got name %q, wanted %q", entity.Name, "AutoTesting")
	}
}

func Test_GormGenericEntityRepository_GetById(t *testing.T) {
	entity := &entities.EthernetSwitch{}
	err := testRepo.GetById(entity, repoTestInsertedId)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if entity.Name != "AutoTesting" {
		t.Errorf("got name %q, wanted %q", entity.Name, "AutoTesting")
	}
}

func Test_GormGenericEntityRepository_Update(t *testing.T) {
	entity := &entities.EthernetSwitch{
		Entity: base.Entity{ID: repoTestInsertedId},
		Name:   "AutoTestingEdited",
	}
	err := testRepo.Update(entity)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
	if entity.Name != "AutoTestingEdited" {
		t.Errorf("got name %q, wanted %q", entity.Name, "AutoTestingEdited")
	}
}

func Test_GormGenericEntityRepository_GetAll(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	err := testRepo.GetAll(entityArr)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
	if len(*entityArr) != 1 {
		t.Errorf("got count %d, wanted %d", len(*entityArr), 1)
	}
}

func Test_GormGenericEntityRepository_GetList(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "", 1, 10, "")
	if count != 1 {
		t.Errorf("got %q element, wanted %q element", count, 1)
	}
	if len(*entityArr) != 1 {
		t.Errorf("got count %d, wanted %d", len(*entityArr), 1)
	}
}

func Test_GormGenericEntityRepository_Delete(t *testing.T) {
	entity := &entities.EthernetSwitch{
		Entity: base.Entity{ID: repoTestInsertedId},
	}
	err := testRepo.Delete(entity)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
}

func Test_GormGenericEntityRepository_GetAllAfterDelete(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	err := testRepo.GetAll(entityArr)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
	if len(*entityArr) != 0 {
		t.Errorf("got count %d, wanted %d", len(*entityArr), 0)
	}
}

func Test_GormGenericEntityRepository_GetListAfterDelete(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "", 1, 10, "")
	if count != 0 {
		t.Errorf("got %q element, wanted %q element", count, 0)
	}
	if len(*entityArr) != 0 {
		t.Errorf("got count %d, wanted %d", len(*entityArr), 0)
	}
}

func Test_GormGenericEntityRepository_Insert20(t *testing.T) {
	for i := 2; i < 22; i++ {
		entity := &entities.EthernetSwitch{
			Entity:      base.Entity{},
			Name:        fmt.Sprintf("AutoTesting %d", i),
			Serial:      "AutoTesting",
			SwitchModel: 0,
			Address:     "10.10.88.2",
			Username:    "AutoTesting",
			Password:    "AutoTesting",
		}
		repoTestInsertedId, _ = testRepo.Insert(entity)
	}
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "", 1, 20, "")
	if count != 20 {
		t.Errorf("got %q element, wanted %q element", count, 20)
	}
}

func Test_GormGenericEntityRepository_GetList_OrderDirection(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "desc", 1, 20, "")
	if count != 20 {
		t.Errorf("got %q element, wanted %q element", count, 20)
	}
	if (*entityArr)[0].ID != 20 {
		t.Errorf("got entity with id %d, wanted %d", (*entityArr)[0].ID, 20)
	}
}

func Test_GormGenericEntityRepository_GetList_Pagination(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "", 2, 10, "")
	if len(*entityArr) != 10 {
		t.Errorf("got %d element, wanted %d element", count, 10)
	}
	if (*entityArr)[0].ID < 9 {
		t.Errorf("got first element on page id %d, wanted %q element", (*entityArr)[0].ID, " or bigger")
	}
}

func Test_GormGenericEntityRepository_GetList_Where(t *testing.T) {
	entityArr := &[]*entities.EthernetSwitch{}
	count, _ := testRepo.GetList(entityArr, "id", "", 2, 10, "name = ?", "AutoTesting 14")
	if count != 1 {
		t.Errorf("got %d element, wanted %d element", count, 1)
	}
	if (*entityArr)[0].Name != "AutoTesting 14" {
		t.Errorf("got entity name %q, wanted %q element", (*entityArr)[0].Name, "AutoTesting 14")
	}
}

func Test_GormGenericEntityRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	sqlDb, err := testRepo.Db.DB()
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
	sqlDb.Close()
	err = os.Remove(testRepoFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}
