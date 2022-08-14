package tests

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
	"time"
)

var (
	testerSwitchService *GenericServiceTest[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]
)

func Test_EthernetSwitchService_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchService_test.db"
	//remove old test db file
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
		new(domain.EthernetSwitch),
		new(domain.EthernetSwitchPort),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.EthernetSwitch]
	repo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testGenDb, logger)
	var service interfaces.IGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]
	switchPortRepository = infrastructure.NewGormGenericRepository[domain.EthernetSwitchPort](testGenDb, logger)
	service, err = services.NewEthernetSwitchService(repo, switchPortRepository, logger)
	if err != nil {
		t.Errorf("create new service failed:  %q", err)
	}
	testerSwitchService = NewGenericServiceTest(service, repo, dbFileName)
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
		return
	}
	if errors.As(err, errors.Validation) {
		cont := errors.GetErrorContext(err)
		if _, ok := cont["SwitchModel"]; ok {
			return
		}
		t.Error("awaiting switch model validation error")
	} else {
		t.Error("awaiting validation error")
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
	// here we create a switch port, to make sure it is removed along with the switch
	portCreateDto := domain.EthernetSwitchPort{
		Name:             "AutoTestingPort",
		EthernetSwitchID: testerSwitchService.InsertedID,
		POEType:          "poe",
	}
	switchPortID, err := switchPortRepository.Insert(context.TODO(), portCreateDto)
	if err != nil {
		t.Error(err)
	}

	// this is where the deletion takes place
	err = testerSwitchService.GenericServiceDelete(testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}

	port, err := switchPortRepository.GetByID(context.TODO(), switchPortID)
	if err != nil {
		t.Errorf("failed to receive switch port: %s", err.Error())
	}
	// since we use soft delete we can still get port from repository
	if port == nil {
		t.Errorf("get by id failed: %s", err.Error())
	}
	// for sure that port was successfully deleted we need to compare it DeletedAt field with old date like 1999-01-01
	boundaryDate, err := time.Parse("2006-01-02", "1999-01-01")
	if err != nil {
		t.Errorf("date parse error: %s", err.Error())
	}
	if port.DeletedAt.Before(boundaryDate) {
		t.Error("error, the switch port was not deleted")
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
