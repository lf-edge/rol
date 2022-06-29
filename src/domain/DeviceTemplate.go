package domain

//DeviceTemplate represents yaml device template as a structure
type DeviceTemplate struct {
	//Name template name
	Name string `yaml:"name"`
	//Model device model
	Model string `yaml:"model"`
	//Manufacturer device manufacturer
	Manufacturer string `yaml:"manufacturer"`
	//Description template description
	Description string `yaml:"description"`
	//CPUCount count of cpus
	CPUCount int `yaml:"cpuCount"`
	//CPUModel model of cpu
	CPUModel string `yaml:"cpuModel"`
	//RAM the amount of RAM in GB
	RAM int `yaml:"ram"`
	//NetworkInterfaces slice of device network interfaces
	NetworkInterfaces []DeviceTemplateNetworkInterface `yaml:"networkInterfaces"`
	//Control describes how we control the device
	Control DeviceTemplateControlDesc `yaml:"control"`
	//DiscBootStages slice of boot stage templates for disk boot
	DiscBootStages []BootStageTemplate `yaml:"discBootStages"`
	//NetBootStages slice of boot stage templates for net boot
	NetBootStages []BootStageTemplate `yaml:"netBootStages"`
	//USBBootStages slice of boot stage templates for usb boot
	USBBootStages []BootStageTemplate `yaml:"usbBootStages"`
}

//DeviceTemplateNetworkInterface is a structure that stores information about network interface
type DeviceTemplateNetworkInterface struct {
	//Name of network interface. This field is unique within device template network interfaces
	Name string `yaml:"name"`
	//NetBoot flags whether the interface can be loaded over the network
	NetBoot bool `yaml:"netBoot"`
	//POEIn only one network interface can be mark as POEIn
	POEIn bool `yaml:"poeIn"`
	//Management only one network interface can be mark as management
	Management bool `yaml:"management"`
}

//DeviceTemplateControlDesc is a structure that stores information
//about how to control the device in different situations
type DeviceTemplateControlDesc struct {
	//Emergency how to control device power in case of emergency. As example: POE(For Rpi4), IPMI, ILO or PowerSwitch
	Emergency string `yaml:"emergency"`
	//Power how to control device power. As example: POE(For Rpi4), IPMI, ILO or PowerSwitch
	Power string `yaml:"power"`
	//NextBoot how to change next boot device. As example: IPMI, ILO or NONE.
	//For example, NONE is used for Rpi4, we control next boot by u-boot files in boot stages.
	NextBoot string `yaml:"nextBoot"`
}

//BootStageTemplateFile is a structure that stores the path to the bootstrap file
type BootStageTemplateFile struct {
	//ExistingFileName file name is a real full file path with name on the disk.
	//This path is relative from app directory
	ExistingFileName string `yaml:"existingFileName"`
	//VirtualFileName virtual file name is relative from /<mac-address>/
	VirtualFileName string `yaml:"virtualFileName"`
}

//BootStageTemplate boot stage can be overwritten in runtime by device entity or by device rent entity.
//BootStageTemplate converts to BootStage for device, then we create device entity.
type BootStageTemplate struct {
	//Name of boot stage template
	Name string `yaml:"name"`
	//Description of boot stage template
	Description string `yaml:"description"`
	//Action for this boot stage.
	//Can be: File, CheckPowerSwitch, EmergencyPowerOff,
	//PowerOff, EmergencyPowerOn, PowerOn,
	//CheckManagement
	//
	//For File action:
	//	A stage can only be marked complete if all files have
	//	been downloaded by the device via TFTP or DHCP,
	//	after which the next step can be loaded.
	Action string `yaml:"action"`
	//Files slice of files for boot stage
	Files []BootStageTemplateFile `yaml:"files"`
}
