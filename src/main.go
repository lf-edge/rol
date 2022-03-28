package main

import (
	"rol/app/interfaces/generic"
	"rol/app/services"
	"rol/domain"
	"rol/infrastructure"
	"rol/webapi"
	"rol/webapi/controllers"

	_ "rol/docs"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// RegisterSwitchController registers controller for the switches
// on path /api/v1/switch/
func RegisterSwitchController(controller *controllers.EthernetSwitchGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/switch/", controller.GetList)
	groupRoute.GET("/switch/:id", controller.GetById)
	groupRoute.POST("/switch", controller.Create)
	groupRoute.PUT("/switch/:id", controller.Update)
	groupRoute.DELETE("/switch/:id", controller.Delete)
}

// RegisterHTTPLogController registers controller for getting http logs via api
// on path /api/v1/httplog/
func RegisterHTTPLogController(controller *controllers.HttpLogGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/log/http/", controller.GetList)
	groupRoute.GET("/log/http/:id", controller.GetById)
}

// RegisterAppLogController registers controller for getting application logs via api
// on path /api/v1/applog/
func RegisterAppLogController(controller *controllers.AppLogGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/log/app/", controller.GetList)
	groupRoute.GET("/log/app/:id", controller.GetById)
}

func RegisterLogHooks(logger *logrus.Logger, httpLogRepo generic.IGenericRepository[domain.HttpLog], logRepo generic.IGenericRepository[domain.AppLog], config *domain.AppConfig) {
	if config.Logger.LogsToDatabase {
		httpHook := infrastructure.NewHttpHook(httpLogRepo)
		appHook := infrastructure.NewAppHook(logRepo)
		logger.AddHook(httpHook)
		logger.AddHook(appHook)
	}
}

// StartHttpServer starts a new http server
func StartHttpServer(server *webapi.GinHTTPServer) {
	server.Start()
}

// @title Rack of labs API
// @version version(1.0)
// @description Description of specifications
// @Precautions when using termsOfService specifications

// @contact.name API supporter
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name license(Mandatory)
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1/
func main() {
	app := fx.New(
		fx.Provide(
			// Realizations
			infrastructure.NewYmlConfig,
			infrastructure.NewGormEntityDb,
			infrastructure.NewGormLogDb,
			infrastructure.NewEthernetSwitchRepository,
			infrastructure.NewHttpLogRepository,
			infrastructure.NewAppLogRepository,
			infrastructure.NewLogrusLogger,
			// Application logic
			services.NewEthernetSwitchService,
			services.NewHttpLogService,
			services.NewAppLogService,
			// WEB API -> Server
			webapi.NewGinHTTPServer,
			// WEB API -> Controllers
			controllers.NewEthernetSwitchGinController,
			controllers.NewHTTPLogGinController,
			controllers.NewAppLogGinController,
		),
		fx.Invoke(
			RegisterLogHooks,
			RegisterSwitchController,
			RegisterHTTPLogController,
			RegisterAppLogController,
			StartHttpServer,
		),
	)
	app.Run()
}
