package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//EthernetSwitchRepository repository for EthernetSwitch entity
type EthernetSwitchRepository struct {
	*GormGenericRepository[uuid.UUID, domain.EthernetSwitch]
}

//NewEthernetSwitchRepository constructor for domain.EthernetSwitch GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewEthernetSwitchRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitch] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.EthernetSwitch](db, log)
	return EthernetSwitchRepository{
		genericRepository,
	}
}
