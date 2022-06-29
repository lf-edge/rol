package dtos

//DeviceTemplateDto dto structure for domain.DeviceTemplate
type DeviceTemplateDto struct {
	//Name template name
	Name string
	//Model device model
	Model string
	//Manufacturer device manufacturer
	Manufacturer string
	//Description template description
	Description string
	//CPUCount count of cpus
	CPUCount int
	//CPUModel model of cpu
	CPUModel string
	//RAM the amount of RAM in GB
	RAM int
	//NetworkInterfaces slice of device network interfaces
	NetworkInterfaces []DeviceTemplateNetworkDto
	//Control describes how we control the device
	Control DeviceTemplateControlDto
	//DiscBootStages slice of boot stage templates for disk boot
	DiscBootStages []DeviceTemplateBootStageDto
	//NetBootStages slice of boot stage templates for net boot
	NetBootStages []DeviceTemplateBootStageDto
	//USBBootStages slice of boot stage templates for usb boot
	USBBootStages []DeviceTemplateBootStageDto
}

//DeviceTemplateNetworkDto dto structure for domain.DeviceTemplateNetworkInterface
type DeviceTemplateNetworkDto struct {
	//Name of network interface. This field is unique within device template network interfaces
	Name string
	//NetBoot flags whether the interface can be loaded over the network
	NetBoot bool
	//POEIn only one network interface can be mark as POEIn
	POEIn bool
	//Management only one network interface can be mark as management
	Management bool
}

//DeviceTemplateControlDto dto structure for domain.DeviceTemplateControlDesc
type DeviceTemplateControlDto struct {
	//Emergency how to control device power in case of emergency. As example: POE(For Rpi4), IPMI, ILO or PowerSwitch
	Emergency string
	//Power how to control device power. As example: POE(For Rpi4), IPMI, ILO or PowerSwitch
	Power string
	//NextBoot how to change next boot device. As example: IPMI, ILO or NONE.
	//For example, NONE is used for Rpi4, we control next boot by u-boot files in boot stages.
	NextBoot string
}

//DeviceTemplateBootStageFileDto dto structure for domain.BootStageTemplateFile
type DeviceTemplateBootStageFileDto struct {
	//ExistingFileName file name is a real full file path with name on the disk.
	//This path is relative from app directory
	ExistingFileName string
	//VirtualFileName virtual file name is relative from /<mac-address>/
	VirtualFileName string
}

//DeviceTemplateBootStageDto dto structure for domain.BootStageTemplate
type DeviceTemplateBootStageDto struct {
	//Name of boot stage template
	Name string
	//Description of boot stage template
	Description string
	//Action for this boot stage.
	//Can be: File, CheckPowerSwitch, EmergencyPowerOff,
	//PowerOff, EmergencyPowerOn, PowerOn,
	//CheckManagement
	//
	//For File action:
	//	A stage can only be marked complete if all files have
	//	been downloaded by the device via TFTP or DHCP,
	//	after which the next step can be loaded.
	Action string
	//Files slice of files for boot stage
	Files []DeviceTemplateBootStageFileDto
}
