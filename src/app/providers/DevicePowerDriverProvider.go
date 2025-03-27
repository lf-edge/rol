package providers

import (
	"rol/app/drivers"
	"rol/app/interfaces"
	"rol/domain"
)

type DevicePowerDriverProvider struct {
	poeDriver *drivers.DevicePOEPowerDriver
}

func NewDevicePowerDriverProvider(poeDriver *drivers.DevicePOEPowerDriver) *DevicePowerDriverProvider {
	return &DevicePowerDriverProvider{
		poeDriver: poeDriver,
	}
}

func (p *DevicePowerDriverProvider) GetPowerManagerDriver(device domain.Device) (interfaces.IDevicePowerDriver, error) {
	if device.PowerControlBus == "POE" {
		return p.poeDriver, nil
	}
	return nil, nil
}
