// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormTFTPConfigRepository repository for TFTPServer entity
type GormTFTPConfigRepository struct {
	*GormGenericRepository[uuid.UUID, domain.TFTPConfig]
}

//NewGormTFTPConfigRepository constructor for domain.TFTPServer GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.TFTPServer] - new tftp server repository
func NewGormTFTPConfigRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.TFTPConfig] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.TFTPConfig](db, log)
	return GormTFTPConfigRepository{
		genericRepository,
	}
}
