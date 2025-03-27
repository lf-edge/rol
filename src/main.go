// Program Rack of Labs
package main

import (
	"os"
	"path/filepath"
	"rol/app/drivers"
	"rol/app/providers"
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

// @host 192.168.8.118:8080
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
			infrastructure.NewGormEthernetSwitchRepository,
			infrastructure.NewGormHTTPLogRepository,
			infrastructure.NewGormAppLogRepository,
			infrastructure.NewGormTFTPConfigRepository,
			infrastructure.NewGormTFTPPathRatioRepository,
			infrastructure.NewPinTFTPServerFactory,
			infrastructure.NewLogrusLogger,
			infrastructure.NewGormEthernetSwitchPortRepository,
			infrastructure.NewDeviceTemplateStorage,
			infrastructure.NewYamlHostNetworkConfigStorage,
			infrastructure.NewHostNetworkManager,
			infrastructure.NewGormEthernetSwitchVLANRepository,
			infrastructure.NewEthernetSwitchManagerProvider,
			infrastructure.NewGormDHCP4LeaseRepository,
			infrastructure.NewGormDHCP4ConfigRepository,
			infrastructure.NewCoreDHCP4ServerFactory,
			infrastructure.NewGormDeviceRepository,
			infrastructure.NewGormDeviceNetworkInterfaceRepository,
			infrastructure.NewGormDevicePowerStateRepository,
			// APP
			// Device power drivers
			drivers.NewDevicePOEPowerDriver,
			// Device power driver provider
			providers.NewDevicePowerDriverProvider,
			// Application logic
			services.NewEthernetSwitchService,
			services.NewHTTPLogService,
			services.NewAppLogService,
			services.NewDeviceTemplateService,
			services.NewHostNetworkService,
			services.NewDHCP4ServerService,
			services.NewTFTPServerService,
			services.NewDeviceService,
			// WEB
			// API -> GIN Server
			webapi.NewGinHTTPServer,
			// API -> GIN Controllers
			controllers.NewEthernetSwitchGinController,
			controllers.NewHTTPLogGinController,
			controllers.NewAppLogGinController,
			controllers.NewEthernetSwitchPortGinController,
			controllers.NewDeviceTemplateController,
			controllers.NewHostNetworkVlanController,
			controllers.NewHostNetworkBridgeController,
			controllers.NewHostNetworkController,
			controllers.NewEthernetSwitchVLANGinController,
			controllers.NewDHCP4ServerGinController,
			controllers.NewTFTPServerGinController,
			controllers.NewDeviceGinController,
		),
		fx.Invoke(
			//Register logrus hooks
			infrastructure.RegisterLogHooks,
			//Services initialization
			services.EthernetSwitchServiceInit,
			services.DHCP4ServerServiceInit,
			services.TFTPServerServiceInit,
			//GIN Controllers registration
			controllers.RegisterEthernetSwitchController,
			controllers.RegisterHTTPLogController,
			controllers.RegisterAppLogController,
			controllers.RegisterEthernetSwitchPortController,
			controllers.RegisterDeviceTemplateController,
			controllers.RegisterHostNetworkVlanController,
			controllers.RegisterHostNetworkBridgeController,
			controllers.RegisterHostNetworkController,
			controllers.RegisterEthernetSwitchVLANGinController,
			controllers.RegisterDHCP4ServerGinController,
			controllers.RegisterTFTPServerGinController,
			controllers.RegisterDeviceGinController,
			//Start GIN http server
			webapi.StartHTTPServer,
		),
	)
	app.Run()
}
