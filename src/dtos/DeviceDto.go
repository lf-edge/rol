package dtos

import "github.com/google/uuid"

type DeviceDto struct {
	BaseDto[uuid.UUID]
	//Name - device name
	Name string
	//Model - device model
	Model string
	//Manufacturer - device manufacturer
	Manufacturer string
	//PowerControlBus - power control bus for this device
	PowerControlBus string
	//PowerState current power state of the device
	PowerState string
}
