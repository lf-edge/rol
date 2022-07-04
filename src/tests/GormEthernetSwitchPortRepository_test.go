package tests

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"runtime"
	"testing"
)

var testerSwitchPortRepository *GenericRepositoryTest[domain.EthernetSwitchPort]

func Test_EthernetSwitchPortRepository_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchPort_test.db"
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitchPort),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.EthernetSwitchPort]
	repo = infrastructure.NewEthernetSwitchPortRepository(testGenDb, logger)

	testerSwitchPortRepository = NewGenericRepositoryTest(repo, dbFileName)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchPortRepository_Insert(t *testing.T) {
	entity := domain.EthernetSwitchPort{
		Name:             "AutoTestingPort",
		EthernetSwitchID: uuid.New(),
		POEType:          "poe",
	}
	err := testerSwitchPortRepository.GenericRepositoryInsert(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_GetByID(t *testing.T) {
	err := testerSwitchPortRepository.GenericRepositoryGetByID(testerSwitchPortRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_GetList(t *testing.T) {
	err := testerSwitchPortRepository.GenericRepositoryGetList()
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_Delete(t *testing.T) {
	err := testerSwitchPortRepository.GenericRepositoryDelete(testerSwitchPortRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		entity := domain.EthernetSwitchPort{
			Name:             fmt.Sprintf("AutoTesting_%d", i),
			EthernetSwitchID: uuid.New(),
			POEType:          fmt.Sprint(i),
		}
		err := testerSwitchPortRepository.GenericRepositoryInsert(entity)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_EthernetSwitchPortRepository_Pagination(t *testing.T) {
	err := testerSwitchPortRepository.GenericRepositoryPagination(1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_Filter(t *testing.T) {
	queryBuilder := testerSwitchPortRepository.Repository.NewQueryBuilder(testerSwitchPortRepository.Context)
	queryGroupBuilder := testerSwitchPortRepository.Repository.NewQueryBuilder(testerSwitchPortRepository.Context)
	queryBuilder.Where("poeType", "!=", "5").Where("poeType", "!=", "9").
		WhereQuery(queryGroupBuilder.Where("name", "=", "AutoTesting_6").Or("name", "=", "AutoTesting_9"))
	err := testerSwitchPortRepository.GenericRepositoryFilter(queryBuilder)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchPortRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerSwitchPortRepository.GenericRepositoryCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}
