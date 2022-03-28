package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapEthernetSwitchPortCreateDto writes ethernet switch port create dto fields to entity
//Params
//	dto - ethernet switch port create dto
//	entity - dest ethernet switch port entity
func MapEthernetSwitchPortCreateDto(dto dtos.EthernetSwitchPortCreateDto, entity *domain.EthernetSwitchPort) {
	entity.Name = dto.Name
	entity.EthernetSwitchId = dto.EthernetSwitchId
	entity.PoeType = (domain.EthernetSwitchPoePortType)(dto.PoeType)
}

//MapEthernetSwitchPortUpdateDto writes ethernet switch port update dto fields to entity
//Params
//	dto - ethernet switch port update dto
//	entity - dest ethernet switch port entity
func MapEthernetSwitchPortUpdateDto(dto dtos.EthernetSwitchPortUpdateDto, entity *domain.EthernetSwitchPort) {
	entity.Name = dto.Name
	entity.PoeType = (domain.EthernetSwitchPoePortType)(dto.PoeType)
}

//MapEthernetSwitchPortToDto writes ethernet switch port entity to dto
//Params
//	entity - ethernet switch port entity
//	dto - dest ethernet switch port dto
func MapEthernetSwitchPortToDto(entity domain.EthernetSwitchPort, dto *dtos.EthernetSwitchPortDto) {
	dto.Id = entity.ID
	dto.Name = entity.Name
	dto.PoeType = (int)(entity.PoeType)
	dto.EthernetSwitchId = entity.EthernetSwitchId
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}
