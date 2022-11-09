package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormDHCP4LeaseRepository repository for domain.DHCP4Lease entity
type GormDHCP4LeaseRepository struct {
	*GormGenericRepository[uuid.UUID, domain.DHCP4Lease]
}

//NewGormDHCP4LeaseRepository constructor for domain.DHCP4Lease GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.DHCP4Lease] - new ethernet switch repository
func NewGormDHCP4LeaseRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.DHCP4Lease] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.DHCP4Lease](db, log)
	return &GormDHCP4LeaseRepository{
		genericRepository,
	}
}
