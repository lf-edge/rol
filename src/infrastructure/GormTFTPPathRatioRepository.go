// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormTFTPPathRatioRepository repository for TFTPPathRatio entity
type GormTFTPPathRatioRepository struct {
	*GormGenericRepository[uuid.UUID, domain.TFTPPathRatio]
}

//NewGormTFTPPathRatioRepository constructor for domain.TFTPPathRatio GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.TFTPPathRatio] - new tftp server repository
func NewGormTFTPPathRatioRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.TFTPPathRatio] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.TFTPPathRatio](db, log)
	return GormTFTPPathRatioRepository{
		genericRepository,
	}
}
