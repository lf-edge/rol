package infrastructure

import (
	"rol/app/interfaces"
	"rol/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// This needed only for FX DI framework

//NewEthernetSwitchPortRepository preparing domain.EthernetSwitchPort repository for passing it in DI
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitchPort] - new ethernet switch repository
func NewEthernetSwitchPortRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitchPort] {
	return NewGormGenericRepository[domain.EthernetSwitchPort](db, log)
}
