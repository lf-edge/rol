package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces"
	"rol/domain"
)

//GormAppLogRepository repository for AppLog entity
type GormAppLogRepository struct {
	*GormGenericRepository[uuid.UUID, domain.AppLog]
}

//NewGormAppLogRepository constructor for domain.AppLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.AppLog] - new app log repository
func NewGormAppLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.AppLog] {
	db := dbShell.GetDb()
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.AppLog](db, log)
	return GormAppLogRepository{
		genericRepository,
	}
}
