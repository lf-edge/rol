package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormEthernetSwitchPortRepository repository for EthernetSwitchPort entity
type GormEthernetSwitchPortRepository struct {
	*GormGenericRepository[uuid.UUID, domain.EthernetSwitchPort]
}

//NewGormEthernetSwitchPortRepository constructor for domain.EthernetSwitch GORM generic repository
//
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewGormEthernetSwitchPortRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchPort] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.EthernetSwitchPort](db, log)
	return GormEthernetSwitchPortRepository{
		genericRepository,
	}
}
