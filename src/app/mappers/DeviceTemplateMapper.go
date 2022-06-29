package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapDeviceTemplateToDto writes device template fields in the dto
func MapDeviceTemplateToDto(template domain.DeviceTemplate, dto *dtos.DeviceTemplateDto) {
	dto.Name = template.Name
	dto.Model = template.Model
	dto.Manufacturer = template.Manufacturer
	dto.Description = template.Description
	dto.CPUCount = template.CPUCount
	dto.CPUModel = template.CPUModel
	dto.RAM = template.RAM
	dto.NetworkInterfaces = mapDeviceTemplateNetworkInterface(template.NetworkInterfaces)
	dto.Control.Power = template.Control.Power
	dto.Control.NextBoot = template.Control.NextBoot
	dto.Control.Emergency = template.Control.Emergency
	dto.DiscBootStages = mapDeviceTemplateBootStage(template.DiscBootStages)
	dto.NetBootStages = mapDeviceTemplateBootStage(template.NetBootStages)
	dto.USBBootStages = mapDeviceTemplateBootStage(template.USBBootStages)
}

func mapDeviceTemplateNetworkInterface(netInterfaces []domain.DeviceTemplateNetworkInterface) []dtos.DeviceTemplateNetworkDto {
	out := make([]dtos.DeviceTemplateNetworkDto, len(netInterfaces))
	for i, netInter := range netInterfaces {
		out[i].Name = netInter.Name
		out[i].NetBoot = netInter.NetBoot
		out[i].POEIn = netInter.POEIn
		out[i].Management = netInter.Management
	}
	return out
}

func mapDeviceTemplateBootStage(bootStages []domain.BootStageTemplate) []dtos.DeviceTemplateBootStageDto {
	out := make([]dtos.DeviceTemplateBootStageDto, len(bootStages))
	for i, stage := range bootStages {
		out[i].Name = stage.Name
		out[i].Description = stage.Description
		out[i].Action = stage.Action
		out[i].Files = mapDeviceTemplateBootStageFile(stage.Files)
	}
	return out
}

func mapDeviceTemplateBootStageFile(bootStageFiles []domain.BootStageTemplateFile) []dtos.DeviceTemplateBootStageFileDto {
	out := make([]dtos.DeviceTemplateBootStageFileDto, len(bootStageFiles))
	for i, file := range bootStageFiles {
		out[i].ExistingFileName = file.ExistingFileName
		out[i].VirtualFileName = file.VirtualFileName
	}
	return out
}
