package tests

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"rol/app/services"
	"rol/domain"
	"rol/infrastructure"
	"strings"
	"testing"
)

var (
	serviceTemplatesCount int
	deviceTemplateService *services.DeviceTemplateService
)

func Test_DeviceTemplateService_Prepare(t *testing.T) {
	serviceTemplatesCount = 30
	err := createXDeviceTemplatesForTest(serviceTemplatesCount)
	if err != nil {
		t.Errorf("creating templates failed: %s", err)
	}
	log := logrus.New()
	storage, err = infrastructure.NewYamlGenericTemplateStorage[domain.DeviceTemplate]("devices", log)
	if err != nil {
		t.Errorf("creating templates storage failed: %s", err.Error())
	}
	deviceTemplateService, err = services.NewDeviceTemplateService(storage, log)
	if err != nil {
		t.Errorf("creating templates service failed: %s", err.Error())
	}
}

func Test_DeviceTemplateService_GetByName(t *testing.T) {
	fileName := fmt.Sprintf("AutoTesting_%d", serviceTemplatesCount/2)
	nameSlice := strings.Split(fileName, ".")
	name := nameSlice[0]
	template, err := deviceTemplateService.GetByName(context.TODO(), fileName)
	if err != nil {
		t.Errorf("get by name failed: %s", err)
	}
	obtainedName := reflect.ValueOf(*template).FieldByName("Name").String()
	if obtainedName != name {
		t.Errorf("unexpected name %s, expect %s", obtainedName, name)
	}
}

func Test_DeviceTemplateService_GetList(t *testing.T) {
	templates, err := deviceTemplateService.GetList(nil, "", "", "", 1, serviceTemplatesCount)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	if templates == nil {
		t.Error("failed get paginated list")
		return
	}
	if templates.Items == nil {
		t.Error("templates not found")
	}
	if len(*templates.Items) != serviceTemplatesCount {
		t.Errorf("unexpected templates count: %d, expect %d", len(*templates.Items), serviceTemplatesCount)
	}
}

func Test_DeviceTemplateService_Search(t *testing.T) {
	templates, err := deviceTemplateService.GetList(nil, "ValueForSearch", "", "", 1, serviceTemplatesCount)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	if templates.Items == nil {
		t.Error("templates not found")
	}
	if len(*templates.Items) != 1 {
		t.Error("search failed")
	}
}

func Test_DeviceTemplateService_DeleteTemplates(t *testing.T) {
	err := removeAllCreatedDeviceTestTemplates()
	if err != nil {
		t.Errorf("deleting device templates failed: %s", err)
	}
}
