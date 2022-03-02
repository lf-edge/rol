package main

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"os"
	"rol/app/interfaces/generic"
	"rol/app/services"
	"rol/infrastructure"
	"rol/webapi"
)

func main() {
	// We need use service as interface, and not as the struct, then we can see implementation errors.
	var service generic.IGenericEntityService = nil
	var repository generic.IGenericEntityRepository = nil
	// Setup sql connection
	dsn := "root:67Edh68Tyt69@tcp(localhost:3306)/rolDb?charset=utf8mb4&parseTime=True&loc=Local"
	gormSqlConnection := mysql.Open(dsn)
	// Setup generic repo (infrastructure layer)
	repository, _ = infrastructure.NewGormGenericEntityRepository(gormSqlConnection)
	//Setup Generic service (business layer, i.e. app)
	service, _ = services.NewGenericEntityService(repository)
	// Setup logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	// Setup http server
	httpServer := webapi.NewHttpServer(logger, &service)
	httpServer.Start()
}
