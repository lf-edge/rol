package tests

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"reflect"
	"rol/app/errors"
	"rol/infrastructure"
	"runtime"
)

//GormSQLiteGenericRepositoryTester tester for gorm sqlite repository
type GormSQLiteGenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]] struct {
	*GenericRepositoryTester[IDType, EntityType]
	dbFileName string
}

//Dispose all tester stuff
func (t *GormSQLiteGenericRepositoryTester[IDType, EntityType]) Dispose() error {
	err := t.GenericRepositoryTester.Dispose()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to remove dispose generic repository tester")
	}
	//remove old test db file
	_, filename, _, _ := runtime.Caller(1)
	if _, err = os.Stat(path.Join(path.Dir(filename), t.dbFileName)); err == nil {
		err = os.Remove(t.dbFileName)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to remove old database")
		}
	}
	return nil
}

//NewGormSQLiteGenericRepositoryTester constructor for gorm sqlite repository tester
func NewGormSQLiteGenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]]() *GormSQLiteGenericRepositoryTester[IDType, EntityType] {
	dbFileName := fmt.Sprintf("gorm_generic_%s_repo_test.db", reflect.TypeOf(new(EntityType)).Elem().Name())
	//remove old test db file
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); err == nil {
		err = os.Remove(dbFileName)
		if err != nil {
			panic("failed to remove old database")
		}
	}
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	err = testGenDb.AutoMigrate(
		new(EntityType),
	)
	if err != nil {
		panic(err)
	}
	repo := infrastructure.NewGormGenericRepository[IDType, EntityType](testGenDb, logrus.New())
	genericTester, _ := NewGenericRepositoryTester[IDType, EntityType](repo)
	genericTester.Implementation = "GORM/SQLITE"
	tester := &GormSQLiteGenericRepositoryTester[IDType, EntityType]{
		GenericRepositoryTester: genericTester,
		dbFileName:              dbFileName,
	}
	return tester
}
