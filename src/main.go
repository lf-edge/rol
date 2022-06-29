package main

import (
	"rol/app/services"
	"rol/infrastructure"
	"rol/webapi"
	"rol/webapi/controllers"

	_ "rol/docs"

	"go.uber.org/fx"
)

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
			infrastructure.NewHTTPLogRepository,
			infrastructure.NewAppLogRepository,
			infrastructure.NewLogrusLogger,
			infrastructure.NewEthernetSwitchPortRepository,
			infrastructure.NewDeviceTemplateStorage,
			// Application logic
			services.NewEthernetSwitchService,
			services.NewHTTPLogService,
			services.NewAppLogService,
			services.NewEthernetSwitchPortService,
			services.NewDeviceTemplateService,
			// WEB API -> Server
			webapi.NewGinHTTPServer,
			// WEB API -> Controllers
			controllers.NewEthernetSwitchGinController,
			controllers.NewHTTPLogGinController,
			controllers.NewAppLogGinController,
			controllers.NewEthernetSwitchPortGinController,
		),
		fx.Invoke(
			infrastructure.RegisterLogHooks,
			controllers.RegisterEthernetSwitchController,
			controllers.RegisterHTTPLogController,
			controllers.RegisterAppLogController,
			controllers.RegisterEthernetSwitchPortController,
			webapi.StartHTTPServer,
		),
	)
	app.Run()
}
