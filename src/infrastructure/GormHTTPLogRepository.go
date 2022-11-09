package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces"
	"rol/domain"
)

//GormHTTPLogRepository repository for HTTPLog entity
type GormHTTPLogRepository struct {
	*GormGenericRepository[uuid.UUID, domain.HTTPLog]
}

//NewGormHTTPLogRepository constructor for domain.HTTPLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.HTTPLog] - new http log repository
func NewGormHTTPLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.HTTPLog] {
	db := dbShell.GetDb()
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.HTTPLog](db, log)
	return GormHTTPLogRepository{
		genericRepository,
	}
}
