package infrastructure

import (
	"rol/app/interfaces/generic"
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
func NewEthernetSwitchRepository(db *gorm.DB, log *logrus.Logger) generic.IGenericRepository[domain.EthernetSwitch] {
	return NewGormGenericRepository[domain.EthernetSwitch](db, log)
}

//NewHttpLogRepository constructor for domain.HttpLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.HttpLog] - new http log repository
func NewHttpLogRepository(dbShell *GormFxShell, log *logrus.Logger) generic.IGenericRepository[domain.HttpLog] {
	db := dbShell.GetDb()
	return NewGormGenericRepository[domain.HttpLog](db, log)
}

//NewAppLogRepository constructor for domain.AppLog GORM generic repository
//Params
//	dbShell - gorm database shell
//	log - logrus logger
//Return
//	generic.IGenericRepository[domain.AppLog] - new app log repository
func NewAppLogRepository(dbShell *GormFxShell, log *logrus.Logger) generic.IGenericRepository[domain.AppLog] {
	db := dbShell.GetDb()
	return NewGormGenericRepository[domain.AppLog](db, log)
}
