package infrastructure

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//EthernetSwitchVLANRepository repository for EthernetSwitchPort entity
type EthernetSwitchVLANRepository struct {
	*GormGenericRepository[domain.EthernetSwitchVLAN]
}

//NewEthernetSwitchVLANRepository constructor for domain.EthernetSwitchVLAN GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewEthernetSwitchVLANRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitchVLAN] {
	genericRepository := NewGormGenericRepository[domain.EthernetSwitchVLAN](db, log)
	return EthernetSwitchVLANRepository{
		genericRepository,
	}
}
