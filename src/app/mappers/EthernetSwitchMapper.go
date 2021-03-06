package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapEthernetSwitchUpdateDto writes ethernet switch update dto fields to entity
//Params
//	dto - ethernet switch update dto
//	entity - dest ethernet switch entity
func MapEthernetSwitchUpdateDto(dto dtos.EthernetSwitchUpdateDto, entity *domain.EthernetSwitch) {
	entity.SwitchModel = dto.SwitchModel
	entity.Name = dto.Name
	entity.Address = dto.Address
	//  pragma: allowlist nextline secret
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

//MapEthernetSwitchCreateDto writes ethernet switch create dto fields to entity
//Params
//	dto - ethernet switch create dto
//	entity - dest ethernet switch entity
func MapEthernetSwitchCreateDto(dto dtos.EthernetSwitchCreateDto, entity *domain.EthernetSwitch) {
	entity.SwitchModel = dto.SwitchModel
	entity.Name = dto.Name
	entity.Address = dto.Address
	//  pragma: allowlist nextline secret
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

//MapEthernetSwitchToDto writes ethernet switch entity fields to dto
//Params
//	entity - ethernet switch entity
//	dto - dest ethernet switch dto
func MapEthernetSwitchToDto(entity domain.EthernetSwitch, dto *dtos.EthernetSwitchDto) {
	dto.ID = entity.ID
	dto.Name = entity.Name
	dto.Address = entity.Address
	dto.Username = entity.Username
	dto.SwitchModel = entity.SwitchModel
	dto.Serial = entity.Serial
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}

//MapEthernetSwitchModelToDto writes ethernet switch model fields to dto
//Params
//	entity - ethernet switch model
//	dto - dest ethernet switch dto
func MapEthernetSwitchModelToDto(entity domain.EthernetSwitchModel, dto *dtos.EthernetSwitchModelDto) {
	dto.Code = entity.Code
	dto.Manufacturer = entity.Manufacturer
	dto.Model = entity.Model
}
