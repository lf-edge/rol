package domain

import "github.com/google/uuid"

//DevPowerState power state for Device
type DevPowerState uint

const (
	DevInPowerOnState = DevPowerState(iota)
	DevInPowerOffState
	DevInPowerUnknownState
)

//String convert state to string
func (s DevPowerState) String() string {
	switch s {
	case DevInPowerOnState:
		return "On"
	case DevInPowerOffState:
		return "Off"
	}
	return "unknown"
}

type DevicePowerState struct {
	EntityUUID
	DeviceID uuid.UUID
	DevPowerState
}
