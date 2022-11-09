package infrastructure

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/interfaces"
	"rol/domain"
)

//GormDHCP4ConfigRepository repository for domain.DHCP4Config entity
type GormDHCP4ConfigRepository struct {
	*GormGenericRepository[uuid.UUID, domain.DHCP4Config]
}

//NewGormDHCP4ConfigRepository constructor for domain.DHCP4Config GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.DHCP4Config] - new ethernet switch repository
func NewGormDHCP4ConfigRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.DHCP4Config] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.DHCP4Config](db, log)
	return &GormDHCP4ConfigRepository{
		genericRepository,
	}
}
