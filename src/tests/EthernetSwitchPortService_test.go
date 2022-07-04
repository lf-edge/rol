package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	"testing"
)

var (
	switchPortService    interfaces.IGenericService[dtos.EthernetSwitchPortDto, dtos.EthernetSwitchPortCreateDto, dtos.EthernetSwitchPortUpdateDto, domain.EthernetSwitchPort]
	switchService        interfaces.IGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]
	switchPortRepository interfaces.IGenericRepository[domain.EthernetSwitchPort]
	portID               uuid.UUID
	ethernetSwitchID     uuid.UUID
)

func Test_EthernetSwitchPortService_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchPortService_test.db"
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitchPort),
		new(domain.EthernetSwitch),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()

	switchPortRepository = infrastructure.NewGormGenericRepository[domain.EthernetSwitchPort](testGenDb, logger)
	switchRepo := infrastructure.NewEthernetSwitchRepository(testGenDb, logger)
	switchPortService, err = services.NewEthernetSwitchPortService(switchPortRepository, switchRepo, logger)
	if err != nil {
		t.Errorf("create new switch port service failed:  %q", err)
	}
	switchService, err = services.NewEthernetSwitchService(switchRepo, logger)
	if err != nil {
		t.Errorf("create new switch service failed:  %q", err)
	}

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchPortService_CreatePortWithoutSwitch(t *testing.T) {
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	service := switchPortService.(*services.EthernetSwitchPortService)
	_, err := service.CreatePort(context.TODO(), uuid.New(), dto)
	if err == nil {
		t.Errorf("nil error, expected: failed to get ethernet switch: switch not found")
	}
}

func Test_EthernetSwitchPortService_CreateSwitchForTests(t *testing.T) {
	switchCreateDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "AutoTesting",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "192.111.111.111",
			Username:    "AutoTesting",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoTesting",
	}
	var err error
	ethernetSwitchID, err = switchService.Create(context.TODO(), switchCreateDto)
	if err != nil {
		t.Errorf("create switch failed: %s", err)
	}
}

func Test_EthernetSwitchPortService_CreatePort(t *testing.T) {
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	service := switchPortService.(*services.EthernetSwitchPortService)
	var err error
	portID, err = service.CreatePort(context.TODO(), ethernetSwitchID, dto)
	if err != nil {
		t.Errorf("create port failed: %s", err)
	}
}

func Test_EthernetSwitchPortService_CreateFailByNonUniqueName(t *testing.T) {
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort",
	}}
	service := switchPortService.(*services.EthernetSwitchPortService)
	_, err := service.CreatePort(context.TODO(), ethernetSwitchID, dto)
	if err == nil {
		t.Errorf("nil error, expected: name uniqueness check error")
	}
}

func Test_EthernetSwitchPortService_CreateFailByBadPOEType(t *testing.T) {
	dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe_",
		Name:    "AutoPort",
	}}
	service := switchPortService.(*services.EthernetSwitchPortService)
	_, err := service.CreatePort(context.TODO(), ethernetSwitchID, dto)
	if err == nil {
		t.Errorf("nil error, expected: dto POEType field validation error")
	}
}

func Test_EthernetSwitchPortService_UpdatePort(t *testing.T) {
	dto := dtos.EthernetSwitchPortUpdateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType: "poe",
		Name:    "AutoPort2.0",
	}}
	service := switchPortService.(*services.EthernetSwitchPortService)
	err := service.UpdatePort(context.TODO(), ethernetSwitchID, portID, dto)
	if err != nil {
		t.Errorf("update port failed: %s", err)
	}

	port, err := switchPortService.GetByID(context.TODO(), portID)
	if err != nil {
		t.Errorf("failed to get port: %s", err)
	}
	if port.Name != "AutoPort2.0" {
		t.Errorf("update port failed: unexpected name, got '%s', expect 'AutoPort2.0'", port.Name)
	}
}

func Test_EthernetSwitchPortService_GetPorts(t *testing.T) {
	service := switchPortService.(*services.EthernetSwitchPortService)
	for i := 1; i < 10; i++ {
		dto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
			POEType: "poe",
			Name:    fmt.Sprintf("AutoPort_%d", i),
		}}
		_, err := service.CreatePort(context.TODO(), ethernetSwitchID, dto)
		if err != nil {
			t.Errorf("create port failed: %s", err)
		}
	}

	ports, err := service.GetPorts(context.TODO(), ethernetSwitchID, "", "", "", 1, 10)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if len(*ports.Items) != 10 {
		t.Errorf("get ports failed: wrong number of items, got %d, expect 10", len(*ports.Items))
	}
}

func Test_EthernetSwitchPortService_Search(t *testing.T) {
	service := switchPortService.(*services.EthernetSwitchPortService)

	ports, err := service.GetPorts(context.TODO(), ethernetSwitchID, "2.0", "", "", 1, 10)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if len(*ports.Items) != 1 {
		t.Errorf("get ports failed: wrong number of items, got %d, expect 1", len(*ports.Items))
	}
	if (*ports.Items)[0].Name != "AutoPort2.0" {
		t.Errorf("get port by ID failed: unexpected name, got '%s', expect 'AutoPort2.0'", (*ports.Items)[0].Name)
	}
}

func Test_EthernetSwitchPortService_GetPortByID(t *testing.T) {
	service := switchPortService.(*services.EthernetSwitchPortService)
	port, err := service.GetPortByID(context.TODO(), ethernetSwitchID, portID)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if port.Name != "AutoPort2.0" {
		t.Errorf("get port by ID failed: unexpected name, got '%s', expect 'AutoPort2.0'", port.Name)
	}
}

func Test_EthernetSwitchPortService_CloseConnectionAndRemoveDb(t *testing.T) {
	if err := switchPortRepository.CloseDb(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove("ethernetSwitchPortService_test.db"); err != nil {
		t.Errorf("remove db failed:  %s", err)
	}
}
