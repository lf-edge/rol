package infrastructure

import (
	"rol/app/interfaces"
	"rol/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// This needed only for FX DI framework

//NewEthernetSwitchRepository constructor for domain.EthernetSwitch GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitch] - new ethernet switch repository
func NewEthernetSwitchRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitch] {
	return NewGormGenericRepository[domain.EthernetSwitch](db, log)
}

//NewHTTPLogRepository constructor for domain.HTTPLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.HTTPLog] - new http log repository
func NewHTTPLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[domain.HTTPLog] {
	db := dbShell.GetDb()
	return NewGormGenericRepository[domain.HTTPLog](db, log)
}

//NewAppLogRepository constructor for domain.AppLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.AppLog] - new app log repository
func NewAppLogRepository(dbShell *GormFxShell, log *logrus.Logger) interfaces.IGenericRepository[domain.AppLog] {
	db := dbShell.GetDb()
	return NewGormGenericRepository[domain.AppLog](db, log)
}

//NewEthernetSwitchPortRepository preparing domain.EthernetSwitchPort repository for passing it in DI
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.EthernetSwitchPort] - new ethernet switch repository
func NewEthernetSwitchPortRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[domain.EthernetSwitchPort] {
	return NewGormGenericRepository[domain.EthernetSwitchPort](db, log)
}
