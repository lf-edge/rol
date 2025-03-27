// Package mappers uses for entity <--> dto conversions
package mappers

import (
	"github.com/google/uuid"
	"rol/domain"
	"rol/dtos"
)

//MapDeviceToDto writes Device entity to dto
//Params
//	entity - Device entity
//	dto - dest Device dto
func MapDeviceToDto(entity domain.Device, dto *dtos.DeviceDto) {
	mapEntityToBaseDto[uuid.UUID](entity, &dto.BaseDto)
	dto.Name = entity.Name
	dto.Model = entity.Model
	dto.PowerControlBus = entity.PowerControlBus
}

//MapDeviceCreateDtoToEntity writes Device create dto fields to entity
//Params
//	dto - Device create dto
//	entity - dest Device entity
func MapDeviceCreateDtoToEntity(dto dtos.DeviceCreateDto, entity *domain.Device) {
	entity.Name = dto.Name
	entity.Model = dto.Model
	entity.PowerControlBus = dto.PowerControlBus
}

//MapDeviceUpdateDtoToEntity writes Device update dto fields to entity
//Params
//	dto - Device update dto
//	entity - dest Device entity
func MapDeviceUpdateDtoToEntity(dto dtos.DeviceUpdateDto, entity *domain.Device) {
	entity.Name = dto.Name
	entity.Model = dto.Model
	entity.PowerControlBus = dto.PowerControlBus
}
