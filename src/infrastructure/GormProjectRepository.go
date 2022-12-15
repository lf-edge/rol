// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormProjectRepository repository for domain.Project entity
type GormProjectRepository struct {
	*GormGenericRepository[uuid.UUID, domain.Project]
}

//NewProjectRepository constructor for domain.Project GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.Project] - new project repository
func NewProjectRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.Project] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.Project](db, log)
	return GormProjectRepository{
		genericRepository,
	}
}
