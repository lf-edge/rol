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
	entity.EthernetSwitchID = dto.EthernetSwitchID
	entity.POEType = dto.POEType
}

//MapEthernetSwitchPortUpdateDto writes ethernet switch port update dto fields to entity
//Params
//	dto - ethernet switch port update dto
//	entity - dest ethernet switch port entity
func MapEthernetSwitchPortUpdateDto(dto dtos.EthernetSwitchPortUpdateDto, entity *domain.EthernetSwitchPort) {
	entity.Name = dto.Name
	entity.POEType = dto.POEType
}

//MapEthernetSwitchPortToDto writes ethernet switch port entity to dto
//Params
//	entity - ethernet switch port entity
//	dto - dest ethernet switch port dto
func MapEthernetSwitchPortToDto(entity domain.EthernetSwitchPort, dto *dtos.EthernetSwitchPortDto) {
	dto.ID = entity.ID
	dto.Name = entity.Name
	dto.POEType = entity.POEType
	dto.EthernetSwitchID = entity.EthernetSwitchID
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
}
