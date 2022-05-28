package infrastructure

import (
	"github.com/sirupsen/logrus"
	"rol/app/interfaces"
	"rol/domain"
)

//AppLogRepository repository for AppLog entity
type AppLogRepository struct {
	*GormGenericRepository[domain.AppLog]
}

//NewAppLogRepository constructor for domain.AppLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.AppLog] - new app log repository
func NewAppLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[domain.AppLog] {
	db := dbShell.GetDb()
	genericRepository := NewGormGenericRepository[domain.AppLog](db, log)
	return AppLogRepository{
		genericRepository,
	}
}
