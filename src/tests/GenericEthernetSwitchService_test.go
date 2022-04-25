package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testServiceRepoFileName string
var testServiceRepoDbConnection gorm.Dialector
var testServiceRepo interfaces.IGenericRepository[domain.EthernetSwitch]
var testService interfaces.IGenericService[
	dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch]
var testServiceCtx context.Context
var insertedId uuid.UUID

func Test_GenericEthernetSwitchService_Prepare(t *testing.T) {
	testServiceRepoFileName = "service_test.db"
	testServiceRepoDbConnection = sqlite.Open(testServiceRepoFileName)
	testDb, err := gorm.Open(testServiceRepoDbConnection, &gorm.Config{})
	err = testDb.AutoMigrate(
		&domain.EthernetSwitch{},
		&domain.EthernetSwitchPort{},
	)
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	testServiceRepo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testDb, logger)
	testService, err = services.NewEthernetSwitchService(testServiceRepo, logger)

	if err != nil {
		t.Errorf("get service failed: %q", err)
	}
	testServiceCtx = context.TODO()

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), testServiceRepoFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(testServiceRepoFileName)
	if err != nil {
		t.Errorf("remove db failed: %q", err)
	}
}

func Test_GenericEthernetSwitchService_Create(t *testing.T) {
	dto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "TestSwitch",
			Serial:      "101",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "Test",
		},
		Password: "Test",
	}
	var err error
	insertedId, err = testService.Create(testServiceCtx, dto)
	if err != nil {
		t.Errorf("create failed: %q", err)
	}
}

func Test_GenericEthernetSwitchService_GetById(t *testing.T) {
	dto, err := testService.GetById(testServiceCtx, insertedId)
	if err != nil {
		t.Errorf("get by id failed: %q", err)
	}
	if dto.Id != insertedId {
		t.Errorf("unexpected entity ID %v, expect 1", dto.Id)
	}
}

func Test_GenericEthernetSwitchService_Update(t *testing.T) {
	updDto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "101",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "Test",
		},
		Password: "Test",
	}
	err := testService.Update(testServiceCtx, updDto, insertedId)
	if err != nil {
		t.Errorf("get by id failed: %q", err)
	}
	dto, err := testService.GetById(testServiceCtx, insertedId)
	if err != nil {
		t.Errorf("get by id failed: %q", err)
	}
	if dto.Name != "AutoTesting" {
		t.Errorf("unexpected entity name %v, expect AutoTesting", dto.Name)
	}
}

func Test_GenericEthernetSwitchService_Delete(t *testing.T) {
	err := testService.Delete(testServiceCtx, insertedId)
	if err != nil {
		t.Errorf("delete failed: %q", err)
	}
	dto, err := testService.GetById(testServiceCtx, insertedId)
	if dto != nil {
		t.Errorf("unexpected entity %v, expect nil", dto)
	}
}

func Test_GenericEthernetSwitchService_Create20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		dto := dtos.EthernetSwitchCreateDto{
			EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
				Name:        fmt.Sprintf("TestSwitch_%d", i),
				Serial:      "101",
				SwitchModel: 0,
				Address:     "123.123.123.123",
				Username:    "Test",
			},
			Password: "Test",
		}
		_, err := testService.Create(testServiceCtx, dto)
		if err != nil {
			t.Errorf("create failed: %q", err)
		}
	}
}

func Test_GenericEthernetSwitchService_GetList(t *testing.T) {
	data, err := testService.GetList(testServiceCtx, "", "CreatedAt", "desc", 1, 10)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if data.Total != 20 {
		t.Errorf("get list failed: total items %d, expect 20", data.Total)
	}
	if data.Page != 1 {
		t.Errorf("get list failed: page %d, expect 1", data.Page)
	}
	if data.Size != 10 {
		t.Errorf("get list failed: size %d, expect 10", data.Size)
	}

	item := (*data.Items)[0]
	if item.Name != "TestSwitch_20" {
		t.Errorf("get list sort failed: get %s, expect TestSwitch_20", item.Name)
	}
}

func Test_GenericEthernetSwitchService_Search(t *testing.T) {
	data, err := testService.GetList(testServiceCtx, "TestSwitch_5", "", "", 1, 10)
	if err != nil {
		t.Errorf("get list failed: %q", err)
	}
	if len(*data.Items) != 1 {
		t.Errorf("get list search failed: size %d, expect 1", len(*data.Items))
	}

	item := (*data.Items)[0]
	if item.Name != "TestSwitch_5" {
		t.Errorf("get list search failed: get %s, expect TestSwitch_5", item.Name)
	}
}

func Test_GenericEthernetSwitchService_CloseConnectionAndRemoveDb(t *testing.T) {
	if err := testServiceRepo.CloseDb(); err != nil {
		t.Errorf("close db failed:  %q", err)
	}
	if err := os.Remove(testServiceRepoFileName); err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}
