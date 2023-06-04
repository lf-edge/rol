package drivers

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//DevicePOEPowerDriver - power control driver for devices that uses POE as power supply
type DevicePOEPowerDriver struct {
	provider        interfaces.IEthernetSwitchManagerProvider
	devEthRepo      interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface]
	switchPortsRepo interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchPort]
	devPowerHistory interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState]
}

//NewDevicePOEPowerDriver - constructor for DevicePOEPowerDriver
func NewDevicePOEPowerDriver(
	provider interfaces.IEthernetSwitchManagerProvider,
	devEthRepo interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface],
	switchPortsRepo interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitchPort],
	devPowerHistory interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState],
) *DevicePOEPowerDriver {
	return &DevicePOEPowerDriver{
		provider:        provider,
		devEthRepo:      devEthRepo,
		switchPortsRepo: switchPortsRepo,
		devPowerHistory: devPowerHistory,
	}
}

func (p *DevicePOEPowerDriver) getDeviceEthPort(ctx context.Context, device domain.Device) (domain.DeviceNetworkInterface, error) {
	devPort := domain.DeviceNetworkInterface{}
	devPorts, err := p.devEthRepo.GetList(ctx, "", "", 1, 1, nil)
	if err != nil {
		return devPort, errors.Internal.Wrapf(err, "failed to get ethernet ports for device: %s", device.ID)
	}
	if len(devPorts) < 1 {
		return devPort, errors.Internal.Wrapf(err, "device with id: %s not have POE ports", device.ID)
	}
	return devPorts[0], nil
}

func (p *DevicePOEPowerDriver) getConnectedSwitchPort(ctx context.Context, devPort domain.DeviceNetworkInterface) (domain.EthernetSwitchPort, error) {
	ethPortQuery := p.switchPortsRepo.NewQueryBuilder(ctx)
	ethPortQuery.Where("EthernetSwitchID", "==", devPort.ConnectedSwitchId)
	switchPort, err := p.switchPortsRepo.GetByIDExtended(ctx, devPort.ConnectedSwitchPortId, ethPortQuery)
	if err != nil {
		return domain.EthernetSwitchPort{}, errors.Internal.Wrapf(err, "failed to get connected switch port for device: %s", devPort.DeviceID)
	}
	return switchPort, nil
}

func (p *DevicePOEPowerDriver) getSwitchPortNameAndManager(ctx context.Context, dev domain.Device) (
	string, interfaces.IEthernetSwitchManager, error,
) {
	devEthPort, err := p.getDeviceEthPort(ctx, dev)
	if err != nil {
		return "", nil, err
	}
	switchEthPort, err := p.getConnectedSwitchPort(ctx, devEthPort)
	if err != nil {
		return "", nil, errors.Internal.Wrapf(err, "get power state of device with id %s failed", dev.ID.String())
	}
	switchManager, err := p.provider.Get(ctx, devEthPort.ConnectedSwitchId)
	if err != nil {
		return "", nil, errors.Internal.Wrapf(err, "get ethernet switch provider failed for device: %s", dev.ID.String())
	}
	return switchEthPort.Name, switchManager, nil
}

//GetPowerState - get device power status from power supply controller
func (p *DevicePOEPowerDriver) GetPowerState(ctx context.Context, dev domain.Device) (domain.DevPowerState, error) {
	devUnkState := domain.DevInPowerUnknownState
	portName, switchManager, err := p.getSwitchPortNameAndManager(ctx, dev)
	if err != nil {
		return devUnkState, err
	}
	poeStatus, err := switchManager.GetPOEPortStatus(portName)
	if err != nil {
		return devUnkState, errors.Internal.Wrapf(err, "failed to get power state of device with id: %s", dev.ID.String())
	}
	switch poeStatus {
	case "enable":
		return domain.DevInPowerOnState, nil
	}
	return domain.DevInPowerOffState, nil
}

//PowerOff device
func (p *DevicePOEPowerDriver) PowerOff(ctx context.Context, dev domain.Device) error {
	portName, switchManager, err := p.getSwitchPortNameAndManager(ctx, dev)
	if err != nil {
		return err
	}
	err = switchManager.DisablePOEPort(portName)
	if err != nil {
		return errors.Internal.Wrapf(err, "disable POE power for device %s failed", dev.ID)
	}
	_, err = p.devPowerHistory.Insert(ctx, domain.DevicePowerState{
		DeviceID:      dev.ID,
		DevPowerState: domain.DevInPowerOffState,
	})
	if err != nil {
		return errors.Internal.Wrapf(err, "save POE power history for device %s failed", dev.ID)
	}
	return nil
}

//PowerOn device
func (p *DevicePOEPowerDriver) PowerOn(ctx context.Context, dev domain.Device) error {
	portName, switchManager, err := p.getSwitchPortNameAndManager(ctx, dev)
	if err != nil {
		return err
	}
	err = switchManager.EnablePOEPort(portName, "poe+")
	if err != nil {
		return errors.Internal.Wrapf(err, "enable POE power for device %s failed", dev.ID)
	}
	_, err = p.devPowerHistory.Insert(ctx, domain.DevicePowerState{
		DeviceID:      dev.ID,
		DevPowerState: domain.DevInPowerOnState,
	})
	if err != nil {
		return errors.Internal.Wrapf(err, "save POE power history for device %s failed", dev.ID)
	}
	return nil
}
