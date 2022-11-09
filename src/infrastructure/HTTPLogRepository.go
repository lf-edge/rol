package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces"
	"rol/domain"
)

//HTTPLogRepository repository for HTTPLog entity
type HTTPLogRepository struct {
	*GormGenericRepository[uuid.UUID, domain.HTTPLog]
}

//NewHTTPLogRepository constructor for domain.HTTPLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.HTTPLog] - new http log repository
func NewHTTPLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.HTTPLog] {
	db := dbShell.GetDb()
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.HTTPLog](db, log)
	return HTTPLogRepository{
		genericRepository,
	}
}
