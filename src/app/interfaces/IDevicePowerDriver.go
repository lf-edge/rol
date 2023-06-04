package interfaces

import (
	"context"
	"rol/domain"
)

type IDevicePowerDriver interface {
	GetPowerState(ctx context.Context, device domain.Device) (domain.DevPowerState, error)
	PowerOn(ctx context.Context, device domain.Device) error
	PowerOff(ctx context.Context, device domain.Device) error
}
