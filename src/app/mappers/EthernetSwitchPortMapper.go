package mappers

import (
	"rol/domain/entities"
	"rol/dtos"
)

func MapEthernetSwitchPortCreateDto(dto *dtos.EthernetSwitchPortCreateDto, model *entities.EthernetSwitchPort) {
	model.Name = dto.Name
	model.EthernetSwitchId = dto.EthernetSwitchId
	model.PoeType = dto.PoeType
}

func MapEthernetSwitchPortUpdateDto(dto *dtos.EthernetSwitchPortUpdateDto, model *entities.EthernetSwitchPort) {
	model.Name = dto.Name
	model.PoeType = dto.PoeType
}

func MapEthernetSwitchPortDto(model *entities.EthernetSwitchPort, dto *dtos.EthernetSwitchPortDto) {
	dto.Id = model.ID
	dto.Name = model.Name
	dto.PoeType = model.PoeType
	dto.EthernetSwitchId = model.EthernetSwitchId
	dto.CreatedAt = model.CreatedAt
	dto.UpdatedAt = model.UpdatedAt
}

func MapEthernetSwitchPortArrayDto(models *[]*entities.EthernetSwitchPort, dtoses *[]*dtos.EthernetSwitchPortDto) {
	for i := range *models {
		dto := &dtos.EthernetSwitchPortDto{}
		MapEthernetSwitchPortDto((*models)[i], dto)
		*dtoses = append(*dtoses, dto)
	}
}
