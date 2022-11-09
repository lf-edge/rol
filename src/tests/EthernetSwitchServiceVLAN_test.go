package tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	customErrors "rol/app/errors"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
)

type tEthSwitchService struct {
	service    *services.EthernetSwitchService
	portRepo   interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchPort]
	vlanRepo   interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchVLAN]
	switchRepo interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitch]
	dbPath     string
	portID     uuid.UUID
	switchID   uuid.UUID
	vlanID     uuid.UUID
}

var ethSwitchServiceTester *tEthSwitchService

func Test_EthernetSwitchServiceVLAN_Prepare(t *testing.T) {
	ethSwitchServiceTester = &tEthSwitchService{}
	ethSwitchServiceTester.dbPath = "ethernetSwitchVlanService_test.db"
	dbConnection := sqlite.Open(ethSwitchServiceTester.dbPath)
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
	switchRepo := infrastructure.NewEthernetSwitchRepository(testGenDb, logger)
	portRepo := infrastructure.NewEthernetSwitchPortRepository(testGenDb, logger)
	vlanRepo := infrastructure.NewEthernetSwitchVLANRepository(testGenDb, logger)
	ethSwitchServiceTester.switchRepo = switchRepo
	ethSwitchServiceTester.portRepo = portRepo
	ethSwitchServiceTester.vlanRepo = vlanRepo

	getter := infrastructure.NewEthernetSwitchManagerProvider(switchRepo)
	service, _ := services.NewEthernetSwitchService(switchRepo, portRepo, vlanRepo, getter)
	ethSwitchServiceTester.service = service
	err = services.EthernetSwitchServiceInit(ethSwitchServiceTester.service)
	if err != nil {
		t.Errorf("init service failed:  %q", err)
	}

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), ethSwitchServiceTester.dbPath)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(ethSwitchServiceTester.dbPath)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchServiceVLAN_CreateRelatedEntities(t *testing.T) {
	ethSwitch := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "TestSwitch",
			Serial:      "serial",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "1.1.1.1",
			Username:    "Test",
		},
		//  pragma: allowlist nextline secret
		Password: "TestTest",
	}
	createdSwitch, err := ethSwitchServiceTester.service.Create(context.Background(), ethSwitch)
	if err != nil {
		t.Errorf("create switch failed:  %q", err)
	}
	switchPortCreateDto := dtos.EthernetSwitchPortCreateDto{EthernetSwitchPortBaseDto: dtos.EthernetSwitchPortBaseDto{
		POEType:    "poe",
		Name:       "Gi",
		POEEnabled: false,
	}}
	createdPort, err := ethSwitchServiceTester.service.CreatePort(context.Background(), createdSwitch.ID, switchPortCreateDto)
	if err != nil {
		t.Errorf("create switch port failed:  %q", err)
	}
	ethSwitchServiceTester.switchID = createdSwitch.ID
	ethSwitchServiceTester.portID = createdPort.ID
}

func Test_EthernetSwitchServiceVLAN_CreateFailByWrongID(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: nil,
			TaggedPorts:   nil,
		},
		VlanID: -1,
	}
	_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, vlanDto)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.Validation) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_CreateFailBySamePortsIDs(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{ethSwitchServiceTester.portID},
			TaggedPorts:   []uuid.UUID{ethSwitchServiceTester.portID},
		},
		VlanID: 2,
	}
	_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, vlanDto)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.Validation) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_CreateFailByNonExistentSwitch(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{ethSwitchServiceTester.portID},
			TaggedPorts:   nil,
		},
		VlanID: 2,
	}
	_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), uuid.New(), vlanDto)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.NotFound) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_CreateFailByNonExistentPort(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{uuid.New()},
			TaggedPorts:   nil,
		},
		VlanID: 2,
	}
	_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, vlanDto)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.Validation) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_CreateOK(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{ethSwitchServiceTester.portID},
			TaggedPorts:   []uuid.UUID{},
		},
		VlanID: 2,
	}
	vlan, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, vlanDto)
	if err != nil {
		t.Errorf("failed to create switch vlan: %q", err)
	}
	ethSwitchServiceTester.vlanID = vlan.ID
}

func Test_EthernetSwitchServiceVLAN_CreateFailByVLANIDUniqueness(t *testing.T) {
	vlanDto := dtos.EthernetSwitchVLANCreateDto{
		EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{ethSwitchServiceTester.portID},
			TaggedPorts:   nil,
		},
		VlanID: 2,
	}
	_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, vlanDto)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.Validation) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_GetByIDFailByNonExistentSwitch(t *testing.T) {
	_, err := ethSwitchServiceTester.service.GetVLANByID(context.Background(), uuid.New(), ethSwitchServiceTester.vlanID)
	if err == nil {
		t.Error("nil error acquired")
	}
	if !customErrors.As(err, customErrors.NotFound) {
		t.Error("wrong error type acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_GetByIDOK(t *testing.T) {
	vlan, err := ethSwitchServiceTester.service.GetVLANByID(context.Background(), ethSwitchServiceTester.switchID, ethSwitchServiceTester.vlanID)
	if err != nil {
		t.Error("get vlan by id failed")
		return
	}
	if vlan.UntaggedPorts[0] != ethSwitchServiceTester.portID {
		t.Error("wrong untagged port acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_Update(t *testing.T) {
	updDto := dtos.EthernetSwitchVLANUpdateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
		UntaggedPorts: nil,
		TaggedPorts:   []uuid.UUID{ethSwitchServiceTester.portID},
	}}
	vlan, err := ethSwitchServiceTester.service.UpdateVLAN(context.Background(), ethSwitchServiceTester.switchID, ethSwitchServiceTester.vlanID, updDto)
	if err != nil {
		t.Errorf("failed to update switch vlan: %q", err)
		return
	}
	if len(vlan.UntaggedPorts) != 0 || vlan.TaggedPorts[0] != ethSwitchServiceTester.portID {
		t.Error("vlan update failed: wrong ports acquired")
	}
}

func Test_EthernetSwitchServiceVLAN_Delete(t *testing.T) {
	err := ethSwitchServiceTester.service.DeleteVLAN(context.Background(), ethSwitchServiceTester.switchID, ethSwitchServiceTester.vlanID)
	if err != nil {
		t.Errorf("failed to delete switch vlan: %q", err)
	}
	_, err = ethSwitchServiceTester.service.GetVLANByID(context.Background(), ethSwitchServiceTester.switchID, ethSwitchServiceTester.vlanID)
	if err == nil {
		t.Error("successfully get removed vlan")
	}
}

func Test_EthernetSwitchServiceVLAN_Create20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		dto := dtos.EthernetSwitchVLANCreateDto{
			EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
				UntaggedPorts: []uuid.UUID{ethSwitchServiceTester.portID},
				TaggedPorts:   nil,
			},
			VlanID: i,
		}
		_, err := ethSwitchServiceTester.service.CreateVLAN(context.Background(), ethSwitchServiceTester.switchID, dto)
		if err != nil {
			t.Errorf("failed to create switch vlan: %q", err)
		}
	}
}

func Test_EthernetSwitchServiceVLAN_GetList(t *testing.T) {
	vlans, err := ethSwitchServiceTester.service.GetVLANs(context.Background(), ethSwitchServiceTester.switchID, "", "VlanID", "desc", 1, 20)
	if err != nil {
		t.Errorf("failed to get switch vlan list: %q", err)
		return
	}
	if len(vlans.Items) != 20 {
		t.Error("wrong vlan's count")
		return
	}
	if vlans.Pagination.Size != 20 {
		t.Error("unexpected page size")
		return
	}
	if vlans.Items[0].VlanID != 20 {
		t.Error("order by failed")
	}
}

func Test_EthernetSwitchServiceVLAN_RemoveDb(t *testing.T) {
	if err := ethSwitchServiceTester.switchRepo.Dispose(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := ethSwitchServiceTester.portRepo.Dispose(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := ethSwitchServiceTester.vlanRepo.Dispose(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove(ethSwitchServiceTester.dbPath); err != nil {
		t.Errorf("remove db failed:  %s", err)
	}
}
