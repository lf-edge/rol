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
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"strings"
	"testing"
)

var testerSwitchService *GenericServiceTest[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]

func Test_EthernetSwitchService_Prepare(t *testing.T) {
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
	var service interfaces.IGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]
	service, err = services.NewEthernetSwitchService(repo, logger)
	if err != nil {
		t.Errorf("create new service failed:  %q", err)
	}
	testerSwitchService = NewGenericServiceTest(service, repo, dbFileName)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchService_CreateFailByWrongModel(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "test_serial",
			SwitchModel: "bad_model",
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	id, err := testerSwitchService.GenericServiceCreate(createDto)
	if err == nil {
		err = testerSwitchService.GenericServiceDelete(id)
		if err != nil {
			t.Error(err)
		}
	} else if !strings.Contains(err.Error(), "switch model is not supported") {
		t.Error("awaiting switch model is not supported error")
	}
}

func Test_EthernetSwitchService_CreateOK(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "Auto Testing",
			Serial:      "test_serial",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	id, err := testerSwitchService.GenericServiceCreate(createDto)
	if err != nil {
		t.Error(err)
	} else {
		testerSwitchService.InsertedID = id
	}
}

func Test_EthernetSwitchService_CreateFailByNotUniqueSerial(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "test_serial",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.124",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	id, err := testerSwitchService.GenericServiceCreate(createDto)
	if err == nil {
		secondErr := testerSwitchService.GenericServiceDelete(id)
		if secondErr != nil {
			t.Error(err, secondErr)
		}
		t.Error("created switch with duplicate serial number")
	}
}

func Test_EthernetSwitchService_CreateFailByNotUniqueAddress(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "test_serial1",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	id, err := testerSwitchService.GenericServiceCreate(createDto)
	if err == nil {
		secondErr := testerSwitchService.GenericServiceDelete(id)
		if secondErr != nil {
			t.Error(err, secondErr)
		}
		t.Error("created switch with duplicate address")
	}
}

func Test_EthernetSwitchService_GetByID(t *testing.T) {
	err := testerSwitchService.GenericServiceGetByID(testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Update(t *testing.T) {
	updateDto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTestingUpdated",
			Serial:      "101",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "Test",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	err := testerSwitchService.GenericServiceUpdate(updateDto, testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Delete(t *testing.T) {
	err := testerSwitchService.GenericServiceDelete(testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Create20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		createDto := dtos.EthernetSwitchCreateDto{
			EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
				Name:        fmt.Sprintf("AutoTesting_%d", i),
				Serial:      fmt.Sprintf("auto_serial_%d", i),
				SwitchModel: "unifi_switch_us-24-250w",
				Address:     fmt.Sprintf("123.123.123.%d", i+1),
				Username:    "AutoUser",
			},
			//  pragma: allowlist nextline secret
			Password: "AutoPass",
		}
		_, err := testerSwitchService.GenericServiceCreate(createDto)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_EthernetSwitchService_GetList(t *testing.T) {
	err := testerSwitchService.GenericServiceGetList(20, 1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Search(t *testing.T) {
	err := testerSwitchService.GenericServiceSearch("AutoUser")
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerSwitchService.GenericServiceCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}
