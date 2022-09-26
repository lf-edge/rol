package tests

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
)

var (
	ethSwitchService       *services.EthernetSwitchService
	ethSwitchRepo          interfaces.IGenericRepository[domain.EthernetSwitch]
	ethSwitchPortRepo      interfaces.IGenericRepository[domain.EthernetSwitchPort]
	ethSwitchVlanRepo      interfaces.IGenericRepository[domain.EthernetSwitchVLAN]
	createdEthSwitchPortID uuid.UUID
	ethSwitchID            uuid.UUID
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
		new(domain.EthernetSwitchVLAN),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	ethSwitchRepo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testGenDb, logger)
	ethSwitchPortRepo = infrastructure.NewGormGenericRepository[domain.EthernetSwitchPort](testGenDb, logger)
	ethSwitchVlanRepo = infrastructure.NewGormGenericRepository[domain.EthernetSwitchVLAN](testGenDb, logger)
	if err != nil {
		t.Error("failed to create switch port repository")
	}
	getter := infrastructure.NewEthernetSwitchManagerProvider(ethSwitchRepo)
	ethSwitchService, err = services.NewEthernetSwitchService(ethSwitchRepo, ethSwitchPortRepo, ethSwitchVlanRepo, getter)
	if err != nil {
		t.Errorf("create new service failed:  %q", err)
	}
	err = services.EthernetSwitchServiceInit(ethSwitchService)
	if err != nil {
		t.Errorf("init service failed:  %q", err)
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
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.Create(ctx, createDto)

	if err == nil {
		err = ethSwitchService.Delete(ctx, ethSwitch.ID)
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
		t.Error("expect switch model validation error")
	} else {
		t.Error("expect validation error")
	}
}

func Test_EthernetSwitchService_CreateOK(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "Auto Testing",
			Serial:      "test_serial",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "UserNameSearch",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.Create(ctx, createDto)
	if err != nil {
		t.Error(err)
	} else {
		ethSwitchID = ethSwitch.ID
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
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.Create(ctx, createDto)
	if err == nil {
		secondErr := ethSwitchService.Delete(ctx, ethSwitch.ID)
		if secondErr != nil {
			t.Error(err, secondErr)
		}
		t.Error("created switch with duplicate serial number")
	}
	if errors.As(err, errors.Validation) {
		cont := errors.GetErrorContext(err)
		if _, ok := cont["Serial"]; ok {
			return
		}
		t.Error("expect switch serial validation error")
	} else {
		t.Error("expect validation error")
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
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.Create(ctx, createDto)
	if err == nil {
		secondErr := ethSwitchService.Delete(ctx, ethSwitch.ID)
		if secondErr != nil {
			t.Error(err, secondErr)
		}
		t.Error("created switch with duplicate address")
	}
	if errors.As(err, errors.Validation) {
		cont := errors.GetErrorContext(err)
		if _, ok := cont["Address"]; ok {
			return
		}
		t.Error("expect switch address validation error")
	} else {
		t.Error("expect validation error")
	}
}

func Test_EthernetSwitchService_GetByID(t *testing.T) {
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.GetByID(ctx, ethSwitchID)
	if err != nil {
		t.Error(err)
	}
	if ethSwitch.Name != "Auto Testing" {
		t.Error("unexpected switch name")
	}
}

func Test_EthernetSwitchService_Update(t *testing.T) {
	updateDto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTestingUpdated",
			Serial:      "101",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "UserNameSearch",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	ctx := context.TODO()
	ethSwitch, err := ethSwitchService.Update(ctx, updateDto, ethSwitchID)
	if err != nil {
		t.Error(err)
	}
	if ethSwitch.Name != "AutoTestingUpdated" {
		t.Error("unexpected switch name")
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
		ctx := context.TODO()
		_, err := ethSwitchService.Create(ctx, createDto)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_EthernetSwitchService_GetList(t *testing.T) {
	ctx := context.TODO()
	list, err := ethSwitchService.GetList(ctx, "", "Name", "DESC", 1, 15)
	if err != nil {
		t.Error(err)
	}
	if list.Pagination.Size != 15 {
		t.Error("unexpected page size")
	}
	if len(list.Items) != 15 {
		t.Error("unexpected items count")
	}
	if list.Items[0].Name != "AutoTesting_9" {
		t.Error("list ordering failed")
	}
}

func Test_EthernetSwitchService_Search(t *testing.T) {
	ctx := context.TODO()
	list, err := ethSwitchService.GetList(ctx, "UserNameSearch", "", "", 1, 10)
	if err != nil {
		t.Error(err)
	}
	if len(list.Items) != 1 {
		t.Error("unexpected items count")
	}
	if list.Items[0].Username != "UserNameSearch" {
		t.Error("search failed")
	}
}

func Test_EthernetSwitchService_CreatePortWithoutSwitch(t *testing.T) {
	ctx := context.TODO()
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	_, err := ethSwitchService.CreatePort(ctx, uuid.New(), dto)
	if err == nil {
		t.Errorf("nil error, expected: error when checking the existence of the switch: switch not found")
	}
}

func Test_EthernetSwitchService_CreatePort(t *testing.T) {
	ctx := context.TODO()
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	var err error
	port, err := ethSwitchService.CreatePort(ctx, ethSwitchID, dto)
	if err != nil {
		t.Errorf("create port failed: %s", err)
	}
	createdEthSwitchPortID = port.ID
}

func Test_EthernetSwitchService_CreatePortFailByNonUniqueName(t *testing.T) {
	ctx := context.TODO()
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	_, err := ethSwitchService.CreatePort(ctx, ethSwitchID, dto)
	if err == nil {
		t.Errorf("nil error, expected: name uniqueness check error")
	}
}

func Test_EthernetSwitchService_CreatePortFailByBadPOEType(t *testing.T) {
	ctx := context.TODO()
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe_",
		Name:    "AutoPort",
	}}
	_, err := ethSwitchService.CreatePort(ctx, ethSwitchID, dto)
	if err == nil {
		t.Errorf("nil error, expected: dto POEType field validation error")
	}
}

func Test_EthernetSwitchService_UpdatePort(t *testing.T) {
	ctx := context.TODO()
	dto := dtos.EthernetSwitchPortUpdateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort2.0",
	}}
	_, err := ethSwitchService.UpdatePort(ctx, ethSwitchID, createdEthSwitchPortID, dto)
	if err != nil {
		t.Errorf("update port failed: %s", err)
	}

	port, err := ethSwitchService.GetPortByID(ctx, ethSwitchID, createdEthSwitchPortID)
	if err != nil {
		t.Errorf("failed to get port: %s", err)
	}
	if port.Name != "AutoPort2.0" {
		t.Errorf("update port failed: unexpected name, got '%s', expect 'AutoPort2.0'", port.Name)
	}
}

func Test_EthernetSwitchService_GetPorts(t *testing.T) {
	ctx := context.TODO()
	for i := 1; i < 10; i++ {
		dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
			POEType: "poe",
			Name:    fmt.Sprintf("AutoPort_%d", i),
		}}
		_, err := ethSwitchService.CreatePort(ctx, ethSwitchID, dto)
		if err != nil {
			t.Errorf("create port failed: %s", err)
		}
	}

	ports, err := ethSwitchService.GetPorts(ctx, ethSwitchID, "", "", "", 1, 10)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if len(ports.Items) != 10 {
		t.Errorf("get ports failed: wrong number of items, got %d, expect 10", len(ports.Items))
	}
}

func Test_EthernetSwitchService_SearchPort(t *testing.T) {
	ctx := context.TODO()
	ports, err := ethSwitchService.GetPorts(ctx, ethSwitchID, "Port2.0", "", "", 1, 10)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if len(ports.Items) != 1 {
		t.Errorf("get ports failed: wrong number of items, got %d, expect 1", len(ports.Items))
	}
	if (ports.Items)[0].Name != "AutoPort2.0" {
		t.Errorf("get port by ID failed: unexpected name, got '%s', expect 'AutoPort2.0'", (ports.Items)[0].Name)
	}
}

func Test_EthernetSwitchService_GetPortByID(t *testing.T) {
	ctx := context.TODO()
	port, err := ethSwitchService.GetPortByID(ctx, ethSwitchID, createdEthSwitchPortID)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if port.Name != "AutoPort2.0" {
		t.Errorf("get port by ID failed: unexpected name, got '%s', expect 'AutoPort2.0'", port.Name)
	}
}

func Test_EthernetSwitchService_DeletePort(t *testing.T) {
	ctx := context.TODO()
	err := ethSwitchService.DeletePort(ctx, ethSwitchID, createdEthSwitchPortID)
	if err != nil {
		t.Errorf("port deletion failed : %s", err)
	}
}

func Test_EthernetSwitchService_Delete(t *testing.T) {
	ctx := context.TODO()
	// here we create a switch port, to make sure it is removed along with the switch
	portCreateDto := dtos.EthernetSwitchPortCreateDto{
		EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
			POEType: "poe",
			Name:    "AutoTestingPort",
		},
	}
	switchPort, err := ethSwitchService.CreatePort(ctx, ethSwitchID, portCreateDto)
	if err != nil {
		t.Error(err)
	}

	// this is where the deletion takes place
	err = ethSwitchService.Delete(ctx, ethSwitchID)
	if err != nil {
		t.Error(err)
	}
	//check port deleting
	_, err = ethSwitchService.GetPortByID(ctx, ethSwitchID, switchPort.ID)
	if err == nil {
		t.Error("successfully get removed switch")
	}
}

func Test_EthernetSwitchService_CloseConnectionAndRemoveDb(t *testing.T) {
	if err := ethSwitchRepo.CloseDb(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := ethSwitchPortRepo.CloseDb(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove("ethernetSwitchService_test.db"); err != nil {
		t.Errorf("remove db failed:  %s", err)
	}
}
