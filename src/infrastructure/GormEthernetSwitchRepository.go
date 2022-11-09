package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormEthernetSwitchRepository repository for EthernetSwitch entity
type GormEthernetSwitchRepository struct {
	*GormGenericRepository[uuid.UUID, domain.EthernetSwitch]
}

//NewGormEthernetSwitchRepository constructor for domain.EthernetSwitch GORM generic repository
//
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewGormEthernetSwitchRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitch] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.EthernetSwitch](db, log)
	return GormEthernetSwitchRepository{
		genericRepository,
	}
}
