// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormDeviceRepository repository for domain.Device entity
type GormDeviceRepository struct {
	*GormGenericRepository[uuid.UUID, domain.Device]
}

//NewGormDeviceRepository constructor for domain.Device GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	interfaces.IGenericRepository[uuid.UUID, domain.Device] - new device repository
func NewGormDeviceRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.Device] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.Device](db, log)
	return &GormDeviceRepository{
		genericRepository,
	}
}
