package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormEthernetSwitchVLANRepository repository for EthernetSwitchPort entity
type GormEthernetSwitchVLANRepository struct {
	*GormGenericRepository[uuid.UUID, domain.EthernetSwitchVLAN]
}

//NewGormEthernetSwitchVLANRepository constructor for domain.EthernetSwitchVLAN GORM generic repository
//
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewGormEthernetSwitchVLANRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchVLAN] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.EthernetSwitchVLAN](db, log)
	return GormEthernetSwitchVLANRepository{
		genericRepository,
	}
}
