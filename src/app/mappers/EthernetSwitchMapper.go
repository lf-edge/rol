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
	entity.SwitchModel = (domain.EthernetSwitchModel)(dto.SwitchModel)
	entity.Name = dto.Name
	entity.Address = dto.Address
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

//MapEthernetSwitchCreateDto writes ethernet switch create dto fields to entity
//Params
//	dto - ethernet switch create dto
//	entity - dest ethernet switch entity
func MapEthernetSwitchCreateDto(dto dtos.EthernetSwitchCreateDto, entity *domain.EthernetSwitch) {
	entity.SwitchModel = (domain.EthernetSwitchModel)(dto.SwitchModel)
	entity.Name = dto.Name
	entity.Address = dto.Address
	entity.Password = dto.Password
	entity.Username = dto.Username
	entity.Serial = dto.Serial
}

//MapEthernetSwitchToDto writes ethernet switch entity fields to dto
//Params
//	entity - ethernet switch entity
//	dto - dest ethernet switch dto
func MapEthernetSwitchToDto(entity domain.EthernetSwitch, dto *dtos.EthernetSwitchDto) {
	dto.Id = entity.ID
	dto.Name = entity.Name
	dto.Address = entity.Address
	dto.Username = entity.Username
	dto.SwitchModel = (int)(entity.SwitchModel)
	dto.Serial = entity.Serial
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}
