package infrastructure

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//EthernetSwitchPortRepository repository for EthernetSwitchPort entity
type EthernetSwitchPortRepository struct {
	*GormGenericRepository[domain.EthernetSwitchPort]
}

//NewEthernetSwitchPortRepository constructor for domain.EthernetSwitch GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewEthernetSwitchPortRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitchPort] {
	genericRepository := NewGormGenericRepository[domain.EthernetSwitchPort](db, log)
	return EthernetSwitchPortRepository{
		genericRepository,
	}
}
