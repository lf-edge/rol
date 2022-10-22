package tests

import (
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

var testerDHCP4ConfigRepository *GenericRepositoryTest[domain.DHCP4Config]

func Test_DHCP4ConfigRepository_Prepare(t *testing.T) {
	dbFileName := "dhcp4config_test.db"
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
		new(domain.DHCP4Config),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.DHCP4Config]
	repo = infrastructure.NewGormDHCP4ConfigRepository(testGenDb, logger)

	testerDHCP4ConfigRepository = NewGenericRepositoryTest(repo, dbFileName)
}

func Test_DHCP4ConfigRepository_Insert(t *testing.T) {
	entity := domain.DHCP4Config{
		Interface:  "eth0",
		Gateway:    "10.10.10.1",
		NTP:        "10.10.10.1",
		ServerID:   "10.10.10.1",
		NextServer: "10.10.10.1",
		Mask:       "255.255.255.0",
		DNS:        "8.8.8.8",
		Range:      "10.10.10.2-10.10.10.20",
		Enabled:    false,
		Port:       67,
		LeaseTime:  3600,
	}
	err := testerDHCP4ConfigRepository.GenericRepositoryInsert(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_GetByID(t *testing.T) {
	err := testerDHCP4ConfigRepository.GenericRepositoryGetByID(testerDHCP4ConfigRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_GetList(t *testing.T) {
	err := testerDHCP4ConfigRepository.GenericRepositoryGetList()
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_Delete(t *testing.T) {
	err := testerDHCP4ConfigRepository.GenericRepositoryDelete(testerDHCP4ConfigRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		entity := domain.DHCP4Config{
			Interface:  fmt.Sprintf("AutoTesting_%d", i),
			Gateway:    fmt.Sprintf("10.10.10.%d", i),
			NTP:        "10.10.10.1",
			ServerID:   "10.10.10.1",
			NextServer: "10.10.10.1",
			Mask:       "255.255.255.0",
			DNS:        "8.8.8.8",
			Range:      "10.10.10.2-10.10.10.20",
			Enabled:    false,
			Port:       67,
			LeaseTime:  3600,
		}
		err := testerDHCP4ConfigRepository.GenericRepositoryInsert(entity)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_DHCP4ConfigRepository_Pagination(t *testing.T) {
	err := testerDHCP4ConfigRepository.GenericRepositoryPagination(1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_Filter(t *testing.T) {
	queryBuilder := testerDHCP4ConfigRepository.Repository.NewQueryBuilder(testerDHCP4ConfigRepository.Context)
	queryGroupBuilder := testerDHCP4ConfigRepository.Repository.NewQueryBuilder(testerDHCP4ConfigRepository.Context)
	queryBuilder.Where("Gateway", "!=", "10.10.10.6").Where("LeaseTime", "!=", 2000).
		WhereQuery(queryGroupBuilder.Where("Interface", "=", "AutoTesting_6").Or("Gateway", "=", "10.10.10.9"))
	err := testerDHCP4ConfigRepository.GenericRepositoryFilter(queryBuilder)
	if err != nil {
		t.Error(err)
	}
}

func Test_DHCP4ConfigRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerDHCP4ConfigRepository.GenericRepositoryCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}
