package tests

import (
	"errors"
	"fmt"
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

var testerSwitchRepository *GenericRepositoryTest[domain.EthernetSwitch]

func Test_EthernetSwitchRepository_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchService_test.db"
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitch),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.EthernetSwitch]
	repo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testGenDb, logger)

	testerSwitchRepository = NewGenericRepositoryTest(repo, dbFileName)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchRepository_Insert(t *testing.T) {
	entity := domain.EthernetSwitch{
		Name:        "AutoTesting",
		Serial:      "1",
		SwitchModel: 0,
		Address:     "123.123.123.123",
		Username:    "AutoName",
		Password:    "AutoPass",
		Ports:       nil,
	}
	err := testerSwitchRepository.GenericRepository_Insert(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_GetById(t *testing.T) {
	err := testerSwitchRepository.GenericRepository_GetById(testerSwitchRepository.InsertedId)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_Update(t *testing.T) {
	entity := domain.EthernetSwitch{
		Entity:      domain.Entity{ID: testerSwitchRepository.InsertedId},
		Name:        "AutoTestingUpdated",
		Serial:      "1",
		SwitchModel: 0,
		Address:     "123.123.123.123",
		Username:    "AutoName",
		Password:    "AutoPass",
		Ports:       nil,
	}
	err := testerSwitchRepository.GenericRepository_Update(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_GetList(t *testing.T) {
	err := testerSwitchRepository.GenericRepository_GetList()
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_Delete(t *testing.T) {
	err := testerSwitchRepository.GenericRepository_Delete(testerSwitchRepository.InsertedId)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		entity := domain.EthernetSwitch{
			Name:        fmt.Sprintf("AutoTesting_%d", i),
			Serial:      fmt.Sprint(i),
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "AutoName",
			Password:    "AutoPass",
			Ports:       nil,
		}
		err := testerSwitchRepository.GenericRepository_Insert(entity)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_EthernetSwitchRepository_Pagination(t *testing.T) {
	err := testerSwitchRepository.GenericRepository_Pagination(1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_Filter(t *testing.T) {
	queryBuilder := testerSwitchRepository.Repository.NewQueryBuilder(testerSwitchRepository.Context)
	queryGroupBuilder := testerSwitchRepository.Repository.NewQueryBuilder(testerSwitchRepository.Context)
	queryBuilder.Where("serial", "!=", "5").Where("serial", "!=", "9").
		WhereQuery(queryGroupBuilder.Where("name", "=", "AutoTesting_6").Or("name", "=", "AutoTesting_9"))
	err := testerSwitchRepository.GenericRepository_Filter(queryBuilder)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerSwitchRepository.GenericRepository_CloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}
