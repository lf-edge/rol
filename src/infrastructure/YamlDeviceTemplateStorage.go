package infrastructure

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces"
	"rol/domain"
)

//YamlDeviceTemplateStorage storage for domain.DeviceTemplate
type YamlDeviceTemplateStorage struct {
	interfaces.IGenericTemplateStorage[domain.DeviceTemplate]
}

//NewDeviceTemplateStorage constructor for DeviceTemplateStorage
func NewDeviceTemplateStorage(log *logrus.Logger) (interfaces.IGenericTemplateStorage[domain.DeviceTemplate], error) {
	storage, err := NewYamlGenericTemplateStorage[domain.DeviceTemplate]("devices", log)
	if err != nil {
		return nil, fmt.Errorf("device templates storage creating error: %s", err.Error())
	}
	return &YamlDeviceTemplateStorage{
		storage,
	}, nil
}
