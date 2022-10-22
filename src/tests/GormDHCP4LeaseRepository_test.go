package tests

import (
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
	"time"
)

var testerDHCP4LeaseRepository *GenericRepositoryTest[domain.DHCP4Lease]

func Test_DHCP4LeaseRepository_Prepare(t *testing.T) {
	dbFileName := "DHCP4Lease_test.db"
	//remove old test db file if exist
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); err == nil {
		err = os.Remove(dbFileName)
		if err != nil {
			t.Errorf("remove db failed:  %q", err)
		}
	}
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.DHCP4Lease),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.DHCP4Lease]
	repo = infrastructure.NewGormDHCP4LeaseRepository(testGenDb, logger)

	testerDHCP4LeaseRepository = NewGenericRepositoryTest(repo, dbFileName)
}

func Test_DHCP4LeaseRepository_Insert(t *testing.T) {
	entity := domain.DHCP4Lease{
		IP:            "10.10.10.2",
		MAC:           "FF:00:01:11:22:33",
		Expires:       time.Now(),
		DHCP4ConfigID: uuid.New(),
	}
	err := testerDHCP4LeaseRepository.GenericRepositoryInsert(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_GetByID(t *testing.T) {
	err := testerDHCP4LeaseRepository.GenericRepositoryGetByID(testerDHCP4LeaseRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_GetList(t *testing.T) {
	err := testerDHCP4LeaseRepository.GenericRepositoryGetList()
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_Delete(t *testing.T) {
	err := testerDHCP4LeaseRepository.GenericRepositoryDelete(testerDHCP4LeaseRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		entity := domain.DHCP4Lease{
			IP:            fmt.Sprintf("10.10.10.%d", i),
			MAC:           fmt.Sprintf("FF:00:01:11:22:%d", i), //only for repository tests
			Expires:       time.Now(),
			DHCP4ConfigID: uuid.New(),
		}
		err := testerDHCP4LeaseRepository.GenericRepositoryInsert(entity)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_DHCP4LeaseRepository_Pagination(t *testing.T) {
	err := testerDHCP4LeaseRepository.GenericRepositoryPagination(1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_Filter(t *testing.T) {
	queryBuilder := testerDHCP4LeaseRepository.Repository.NewQueryBuilder(testerDHCP4LeaseRepository.Context)
	queryGroupBuilder := testerDHCP4LeaseRepository.Repository.NewQueryBuilder(testerDHCP4LeaseRepository.Context)
	queryBuilder.Where("IP", "!=", "10.10.10.6").Where("MAC", "!=", "FF:00:01:11:22:7").
		WhereQuery(queryGroupBuilder.Where("IP", "=", "10.10.10.6").Or("IP", "=", "10.10.10.9"))
	err := testerDHCP4LeaseRepository.GenericRepositoryFilter(queryBuilder)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4LeaseRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerDHCP4LeaseRepository.GenericRepositoryCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}
