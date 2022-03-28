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
	err := createDbIfNotExists(cfg.Database.EntityConnectionString, cfg.Database.EntityDbName)
	if err != nil {
		return nil, err
	}
	connectionString := fmt.Sprintf("%s%s%s", cfg.Database.EntityConnectionString, cfg.Database.EntityDbName, cfg.Database.EntityDbParams)
	dialector := mysql.Open(connectionString)
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
}

//NewGormLogDb creates new gorm logs database connection and create tables if necessary
//Params
//	cfg - application configuration
//Return
//	*GormFxShell - gorm database shell
//	error - if an error occurs, otherwise nil
func NewGormLogDb(cfg *domain.AppConfig) (*GormFxShell, error) {
	err := createDbIfNotExists(cfg.Database.LogConnectionString, cfg.Database.LogDbName)
	if err != nil {
		return nil, err
	}
	connectionString := fmt.Sprintf("%s%s%s", cfg.Database.LogConnectionString, cfg.Database.LogDbName, cfg.Database.LogDbParams)
	dialector := mysql.Open(connectionString)
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
	db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	return nil
}
