package infrastructure

import (
	"errors"
	"fmt"
	"rol/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type OurNamingSchema struct {
	schema.NamingStrategy
}

// TableName convert string to table name
func (ns OurNamingSchema) TableName(str string) string {
	defaultNamingStrategy := schema.NamingStrategy{}
	return defaultNamingStrategy.TableName(str)
}

//NewGormEntityDb creates new gorm entity database connection and create tables if necessary
//Params
//	cfg - application configuration
//Return
//	*gorm.DB - gorm database
//	error - if an error occurs, otherwise nil
func NewGormEntityDb(cfg *domain.AppConfig) (*gorm.DB, error) {
	entityCfg := cfg.Database.Entity
	connectionString := fmt.Sprintf("%s:%s@%s(%s:%s)/", entityCfg.Username, entityCfg.Password, entityCfg.Protocol, entityCfg.Hostname, entityCfg.Port)
	err := createDbIfNotExists(connectionString, entityCfg.DbName)
	if err != nil {
		return nil, err
	}
	dialector := mysql.Open(fmt.Sprintf("%s%s%s", connectionString, entityCfg.DbName, entityCfg.Parameters))
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[NewGormEntityDb]: Failed to open db: %v", err))
	}
	err = db.AutoMigrate(
		&domain.EthernetSwitch{},
		&domain.EthernetSwitchPort{},
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[NewGormEntityDb]: Failed to apply db migrations: %v", err))
	}
	return db, nil
	// entityConnectionString: "root:67Edh68Tyt69@tcp(localhost:3306)/"
}

//NewGormLogDb creates new gorm logs database connection and create tables if necessary
//Params
//	cfg - application configuration
//Return
//	*GormFxShell - gorm database shell
//	error - if an error occurs, otherwise nil
func NewGormLogDb(cfg *domain.AppConfig) (*GormFxShell, error) {
	logCfg := cfg.Database.Log
	connectionString := fmt.Sprintf("%s:%s@%s(%s:%s)/", logCfg.Username, logCfg.Password, logCfg.Protocol, logCfg.Hostname, logCfg.Port)
	err := createDbIfNotExists(connectionString, logCfg.DbName)
	if err != nil {
		return nil, err
	}
	dialector := mysql.Open(fmt.Sprintf("%s%s%s", connectionString, logCfg.DbName, logCfg.Parameters))
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[NewGormLogDb]: Failed to open db: %v", err))
	}
	err = db.AutoMigrate(
		&domain.HttpLog{},
		&domain.AppLog{},
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[NewGormLogDb]: Failed to apply db migrations: %v", err))
	}
	return &GormFxShell{dbShell: db}, nil
}

func createDbIfNotExists(connectionString, dbName string) error {
	dialector := mysql.Open(connectionString)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
	if err != nil {
		return err
	}
	return nil
}
