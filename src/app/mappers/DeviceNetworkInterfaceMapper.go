// Package mappers uses for entity <--> dto conversions
package mappers

import (
	"github.com/google/uuid"
	"rol/domain"
	"rol/dtos"
)

//MapDeviceNetworkInterfaceToDto writes DeviceNetworkInterface entity to dto
//Params
//	entity - DeviceNetworkInterface entity
//	dto - dest DeviceNetworkInterface dto
func MapDeviceNetworkInterfaceToDto(entity domain.DeviceNetworkInterface, dto *dtos.DeviceNetworkInterfaceDto) {
	mapEntityToBaseDto[uuid.UUID](entity, &dto.BaseDto)
	dto.Mac = entity.Mac
	dto.ConnectedSwitchPortId = entity.ConnectedSwitchPortId
	dto.ConnectedSwitchId = entity.ConnectedSwitchId
}

//MapDeviceNetworkInterfaceCreateDtoToEntity writes DeviceNetworkInterface create dto fields to entity
//Params
//	dto - DeviceNetworkInterface create dto
//	entity - dest DeviceNetworkInterface entity
func MapDeviceNetworkInterfaceCreateDtoToEntity(dto dtos.DeviceNetworkInterfaceCreateDto, entity *domain.DeviceNetworkInterface) {
	entity.Mac = dto.Mac
	entity.ConnectedSwitchId = dto.ConnectedSwitchId
	entity.ConnectedSwitchPortId = dto.ConnectedSwitchPortId
}

//MapDeviceNetworkInterfaceUpdateDtoToEntity writes DeviceNetworkInterface update dto fields to entity
//Params
//	dto - DeviceNetworkInterface update dto
//	entity - dest DeviceNetworkInterface entity
func MapDeviceNetworkInterfaceUpdateDtoToEntity(dto dtos.DeviceNetworkInterfaceUpdateDto, entity *domain.DeviceNetworkInterface) {
	entity.Mac = dto.Mac
	entity.ConnectedSwitchId = dto.ConnectedSwitchId
	entity.ConnectedSwitchPortId = dto.ConnectedSwitchPortId
}
