// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormDeviceNetworkInterfaceRepository repository for domain.DeviceNetworkInterface entity
type GormDeviceNetworkInterfaceRepository struct {
	*GormGenericRepository[uuid.UUID, domain.DeviceNetworkInterface]
}

//NewGormDeviceNetworkInterfaceRepository constructor for domain.DeviceNetworkInterface GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface] - new device network interface repository
func NewGormDeviceNetworkInterfaceRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.DeviceNetworkInterface](db, log)
	return &GormDeviceNetworkInterfaceRepository{
		genericRepository,
	}
}
