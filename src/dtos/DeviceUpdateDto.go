package dtos

type DeviceUpdateDto struct {
	//Name - device name
	Name string
	//Model - device model
	Model string
	//Manufacturer - device manufacturer
	Manufacturer string
	//PowerControlBus - power control bus for this device
	PowerControlBus string
}
