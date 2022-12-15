// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"os"
	"path"
	"rol/app/errors"
	"rol/domain"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//OurNamingSchema tables, columns naming strategy
type OurNamingSchema struct {
	schema.NamingStrategy
}

// TableName convert string to table name
func (ns OurNamingSchema) TableName(str string) string {
	defaultNamingStrategy := schema.NamingStrategy{}
	return defaultNamingStrategy.TableName(str)
}

func createMySQLDb(cfg domain.MySQL) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@%s(%s:%s)/", cfg.Username, cfg.Password,
		cfg.Protocol, cfg.Hostname, cfg.Port)
	dialector := mysql.Open(fmt.Sprintf("%s%s%s", connectionString, cfg.DbName, cfg.Parameters))
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to open mysql db")
	}
	err = createDbIfNotExists(db, cfg.DbName)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error creating db")
	}
	return db, nil
}

func createSQLiteDb(cfg domain.SQLite) (*gorm.DB, error) {
	dbPath := cfg.Filename
	pathSplit := strings.Split(dbPath, "/")
	dbName := pathSplit[len(pathSplit)-1]

	if dbPath[0:1] != "/" {
		executablePath, _ := os.Executable()
		dbPath = path.Join(path.Dir(executablePath), dbPath)
	}
	err := os.MkdirAll(dbPath[:len(dbPath)-len(dbName)], os.ModePerm)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to create directory")
	}
	dbConnection := sqlite.Open(dbPath)
	db, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		return nil, errors.Internal.Wrap(err, "creating sqlite db failed")
	}
	return db, nil
}

func newGormDb(cfg domain.DbConfig) (*gorm.DB, error) {
	switch cfg.Driver {
	case "mysql":
		return createMySQLDb(cfg.MySQL)
	case "sqlite":
		return createSQLiteDb(cfg.SQLite)
	case "":
		return nil, errors.Internal.New("no db driver specified")
	default:
		return nil, errors.Internal.Newf("unsupported db driver: %s", cfg.Driver)
	}
}

//NewGormEntityDb creates new gorm entity database connection and create tables if necessary
//Params
//	cfg - application configuration
//Return
//	*gorm.DB - gorm database
//	error - if an error occurs, otherwise nil
func NewGormEntityDb(cfg *domain.AppConfig) (*gorm.DB, error) {
	entityCfg := cfg.Database.Entity
	db, err := newGormDb(entityCfg)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to create db")
	}
	err = db.AutoMigrate(
		&domain.TFTPConfig{},
		&domain.TFTPPathRatio{},
		&domain.EthernetSwitch{},
		&domain.EthernetSwitchPort{},
		&domain.EthernetSwitchVLAN{},
		&domain.DHCP4Config{},
		&domain.DHCP4Lease{},
		&domain.Device{},
		&domain.DeviceNetworkInterface{},
		&domain.DevicePowerState{},
	)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to apply db migrations")
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
	logCfg := cfg.Database.Log
	db, err := newGormDb(logCfg)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to create db")
	}
	err = db.AutoMigrate(
		&domain.HTTPLog{},
		&domain.AppLog{},
	)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to apply db migrations")
	}
	return &GormFxShell{dbShell: db}, nil
}

func createDbIfNotExists(db *gorm.DB, dbName string) error {
	err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
	if err != nil {
		return errors.Internal.Wrap(err, "failed execute raw sql")
	}
	return nil
}
