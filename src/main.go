package main

import (
	"os"
	"path/filepath"
	"rol/app/services"
	"rol/domain"
	"rol/infrastructure"
	"rol/webapi"
	"rol/webapi/controllers"

	_ "rol/webapi/swagger"

	"go.uber.org/fx"
)

//GetGlobalDIParameters get global parameters for DI
func GetGlobalDIParameters() domain.GlobalDIParameters {
	filePath, _ := os.Executable()
	return domain.GlobalDIParameters{
		RootPath: filepath.Dir(filePath),
	}
}

// @title Rack of labs API
// @version 0.1.0
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
			// Core
			GetGlobalDIParameters,
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
			infrastructure.NewYamlHostNetworkConfigStorage,
			infrastructure.NewHostNetworkManager,
			// Application logic
			services.NewEthernetSwitchService,
			services.NewHTTPLogService,
			services.NewAppLogService,
			services.NewEthernetSwitchPortService,
			services.NewDeviceTemplateService,
			services.NewHostNetworkVlanService,
			// WEB API -> GIN Server
			webapi.NewGinHTTPServer,
			// WEB API -> GIN Controllers
			controllers.NewEthernetSwitchGinController,
			controllers.NewHTTPLogGinController,
			controllers.NewAppLogGinController,
			controllers.NewEthernetSwitchPortGinController,
			controllers.NewDeviceTemplateController,
			controllers.NewHostNetworkVlanController,
		),
		fx.Invoke(
			//Register logrus hooks
			infrastructure.RegisterLogHooks,
			//Services initialization
			services.EthernetSwitchServiceInit,
			//GIN Controllers registration
			controllers.RegisterEthernetSwitchController,
			controllers.RegisterHTTPLogController,
			controllers.RegisterAppLogController,
			controllers.RegisterEthernetSwitchPortController,
			controllers.RegisterDeviceTemplateController,
			controllers.RegisterHostNetworkVlanController,
			//Start GIN http server
			webapi.StartHTTPServer,
		),
	)
	app.Run()
}
