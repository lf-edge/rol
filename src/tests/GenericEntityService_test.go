package testss

import (
	"gorm.io/driver/sqlite"
	"os"
	"rol/app/services"
	"rol/dtos"
	"rol/infrastructure"
	"testing"
)

var testServiceFileName = "service_test.db"
var testServiceDbConnection = sqlite.Open(testServiceFileName)
var testServiceRepo, _ = infrastructure.NewGormGenericEntityRepository(testServiceDbConnection)
var testService, _ = services.NewGenericEntityService(testServiceRepo)
var serviceTestCreatedId uint = 0

func Test_GenericEntityService_Create(t *testing.T) {
	dto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "TestSwitch",
			Serial:      "TestSerial",
			SwitchModel: 146,
			Address:     "TestAddress",
			Username:    "TestUsername",
		},
		Password: "TestPassword",
	}
	var err error

	serviceTestCreatedId, err = testService.Create(&dto)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if serviceTestCreatedId == 0 {
		t.Errorf("got %q, wanted %q", serviceTestCreatedId, " > 0")
	}
	if dto.Name != "TestSwitch" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestSwitch")
	}
}

func Test_GenericEntityService_GetById(t *testing.T) {
	dto := dtos.EthernetSwitchDto{}

	err := testService.GetById(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if dto.Name != "TestSwitch" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestSwitch")
	}
}

func Test_GenericEntityService_Update(t *testing.T) {
	dto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{Name: "TestEdit"},
	}
	err := testService.Update(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if dto.Name != "TestEdit" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestEdit")
	}
}

func Test_GenericEntityService_GetAll(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	err := testService.GetAll(dtosArr)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if len(*dtosArr) != 1 {
		t.Errorf("got count %d, wanted %d", len(*dtosArr), 1)
	}
}

func Test_GenericEntityService_Delete(t *testing.T) {
	dto := dtos.EthernetSwitchDto{}
	err := testService.Delete(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
}

func Test_GenericEntityService_GetAllAfterDelete(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	err := testService.GetAll(dtosArr)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
	if len(*dtosArr) != 0 {
		t.Errorf("got count %d, wanted %d", len(*dtosArr), 0)
	}
}

func Test_GenericEntityService_CloseConnectionAndRemoveDb(t *testing.T) {
	sqlDb, err := testServiceRepo.Db.DB()
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
	sqlDb.Close()
	err = os.Remove(testServiceFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}
