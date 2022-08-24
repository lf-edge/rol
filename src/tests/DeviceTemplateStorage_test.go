package tests

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"strings"
	"testing"
)

var storageTemplatesCount int
var storage interfaces.IGenericTemplateStorage[domain.DeviceTemplate]

func Test_DeviceTemplateStorage_Prepare(t *testing.T) {
	storageTemplatesCount = 30
	err := createXDeviceTemplatesForTest(storageTemplatesCount)
	if err != nil {
		t.Errorf("creating templates failed: %s", err)
	}
	storage, err = infrastructure.NewYamlGenericTemplateStorage[domain.DeviceTemplate]("devices", logrus.New())
	if err != nil {
		t.Errorf("creating templates storage failed: %s", err.Error())
	}
}

func Test_DeviceTemplateStorage_GetByName(t *testing.T) {
	fileName := fmt.Sprintf("AutoTesting_%d", storageTemplatesCount/2)
	nameSlice := strings.Split(fileName, ".")
	name := nameSlice[0]
	template, err := storage.GetByName(context.TODO(), fileName)
	if err != nil {
		t.Errorf("get by name failed: %s", err)
	}
	obtainedName := reflect.ValueOf(template).FieldByName("Name").String()
	if obtainedName != name {
		t.Errorf("unexpected name %s, expect %s", obtainedName, name)
	}
}

func Test_DeviceTemplateStorage_GetList(t *testing.T) {
	templatesArr, err := storage.GetList(context.TODO(), "", "", 1, storageTemplatesCount, nil)
	if err != nil {
		t.Errorf("get list failed:  %s", err)
	}
	if len(templatesArr) != storageTemplatesCount {
		t.Errorf("array length %d, expect %d", len(templatesArr), storageTemplatesCount)
	}
}

func Test_DeviceTemplateStorage_Pagination(t *testing.T) {
	page := 1
	pageSize := 10
	templatesArrFirstPage, err := storage.GetList(context.TODO(), "CPUCount", "asc", page, pageSize, nil)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	if len(templatesArrFirstPage) != pageSize {
		t.Errorf("array length on %d page %d, expect %d", page, len(templatesArrFirstPage), pageSize)
	}
	templatesArrSecondPage, err := storage.GetList(context.TODO(), "CPUCount", "asc", page+1, pageSize, nil)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	if len(templatesArrSecondPage) != pageSize {
		t.Errorf("array length on next page %d, expect %d", len(templatesArrSecondPage), pageSize)
	}

	firstPageValue := reflect.ValueOf(templatesArrFirstPage).Index(0)
	secondPageValue := reflect.ValueOf(templatesArrSecondPage).Index(0)

	firstPageValueName := firstPageValue.FieldByName("Name").String()
	secondPageValueName := secondPageValue.FieldByName("Name").String()
	if firstPageValueName == secondPageValueName {
		t.Errorf("pagination failed: got same element on second page with Name: %s", firstPageValueName)
	}
	firstPageValueCPU := firstPageValue.FieldByName("CPUCount").Int()
	secondPageValueCPU := secondPageValue.FieldByName("CPUCount").Int()

	if secondPageValueCPU-int64(pageSize) != firstPageValueCPU {
		t.Errorf("pagination failed: unexpected element on second page")
	}
}

func Test_DeviceTemplateStorage_Sort(t *testing.T) {
	templatesArr, err := storage.GetList(context.TODO(), "CPUCount", "asc", 1, storageTemplatesCount, nil)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	if len(templatesArr) != storageTemplatesCount {
		t.Errorf("array length %d, expect %d", len(templatesArr), storageTemplatesCount)
	}
	index := storageTemplatesCount / 2
	name := reflect.ValueOf(templatesArr).Index(index - 1).FieldByName("Name").String()

	if name != fmt.Sprintf("AutoTesting_%d", index) {
		t.Errorf("sort failed: got %s name, expect AutoTesting_%d", name, index)
	}
}

func Test_DeviceTemplateStorage_Filter(t *testing.T) {
	queryBuilder := storage.NewQueryBuilder(context.TODO())
	queryBuilder.
		Where("CPUCount", ">", storageTemplatesCount/2).
		Where("CPUCount", "<", storageTemplatesCount)
	templatesArr, err := storage.GetList(context.TODO(), "", "", 1, storageTemplatesCount, queryBuilder)
	if err != nil {
		t.Errorf("get list failed: %s", err)
	}
	var expectedCount int
	if storageTemplatesCount%2 == 0 {
		expectedCount = storageTemplatesCount/2 - 1
	} else {
		expectedCount = storageTemplatesCount / 2
	}
	if len(templatesArr) != expectedCount {
		t.Errorf("array length %d, expect %d", len(templatesArr), expectedCount)
	}
}

func Test_DeviceTemplateStorage_DeleteTemplates(t *testing.T) {
	err := removeAllCreatedDeviceTestTemplates()
	if err != nil {
		t.Errorf("deleting device templates failed: %s", err)
	}
}

func createXDeviceTemplatesForTest(x int) error {
	executedFilePath, _ := os.Executable()
	templatesDir := path.Join(path.Dir(executedFilePath), "templates", "devices")
	err := os.MkdirAll(templatesDir, 0777)
	if err != nil {
		return fmt.Errorf("creating dir failed: %s", err)
	}
	for i := 1; i <= x; i++ {
		template := domain.DeviceTemplate{
			Name:         fmt.Sprintf("AutoTesting_%d", i),
			Model:        fmt.Sprintf("AutoTesting_%d", i),
			Manufacturer: "Manufacturer",
			Description:  "Description",
			CPUCount:     i,
			CPUModel:     "CPUModel",
			RAM:          i,
			NetworkInterfaces: []domain.DeviceTemplateNetworkInterface{{
				Name:       "Name",
				NetBoot:    false,
				POEIn:      false,
				Management: false,
			}},
			Control: domain.DeviceTemplateControlDesc{
				Emergency: "Emergency",
				Power:     "Power",
				NextBoot:  "NextBoot",
			},
			DiscBootStages: []domain.BootStageTemplate{{
				Name:        "Name",
				Description: "Description",
				Action:      "Action",
				Files: []domain.BootStageTemplateFile{{
					ExistingFileName: "ExistingFileName",
					VirtualFileName:  "VirtualFileName",
				}},
			}},
			NetBootStages: []domain.BootStageTemplate{{
				Name:        "Name",
				Description: "Description",
				Action:      "Action",
				Files: []domain.BootStageTemplateFile{{
					ExistingFileName: "ExistingFileName",
					VirtualFileName:  "VirtualFileName",
				}},
			}},
			USBBootStages: []domain.BootStageTemplate{{
				Name:        "Name",
				Description: "Description",
				Action:      "Action",
				Files: []domain.BootStageTemplateFile{{
					ExistingFileName: "ExistingFileName",
					VirtualFileName:  "VirtualFileName",
				}},
			}},
		}
		if i == 2 {
			template.Description = "ValueForSearch"
		}
		yamlData, err := yaml.Marshal(&template)
		if err != nil {
			return fmt.Errorf("yaml marshal failed: %s", err)
		}
		fileName := path.Join(templatesDir, fmt.Sprintf("AutoTesting_%d.yml", i))
		err = ioutil.WriteFile(fileName, yamlData, 0777)
		if err != nil {
			return fmt.Errorf("create yaml file failed: %s", err)
		}
	}
	return nil
}

func removeAllCreatedDeviceTestTemplates() error {
	executedFilePath, _ := os.Executable()
	templatesDir := path.Join(path.Dir(executedFilePath), "templates", "devices")
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), "AutoTesting_") {
			err := os.Remove(path.Join(templatesDir, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
