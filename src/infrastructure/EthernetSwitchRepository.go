package infrastructure

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//EthernetSwitchRepository repository for EthernetSwitch entity
type EthernetSwitchRepository struct {
	*GormGenericRepository[domain.EthernetSwitch]
}

//NewEthernetSwitchRepository constructor for domain.EthernetSwitch GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewEthernetSwitchRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitch] {
	genericRepository := NewGormGenericRepository[domain.EthernetSwitch](db, log)
	return EthernetSwitchRepository{
		genericRepository,
	}
}
