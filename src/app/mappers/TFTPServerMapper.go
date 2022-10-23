// Package mappers uses for entity <--> dto conversions
package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapTFTPConfigToDto writes TFTP config entity to dto
//
//Params
//	entity - ethernet switch port entity
//	dto - dest ethernet switch port dto
func MapTFTPConfigToDto(entity domain.TFTPConfig, dto *dtos.TFTPServerDto) {
	dto.Address = entity.Address
	dto.Port = entity.Port
	dto.ID = entity.ID
	dto.UpdatedAt = entity.UpdatedAt
	dto.CreatedAt = entity.CreatedAt
	dto.Enabled = entity.Enabled
}

//MapTFTPServerCreateDtoToEntity writes TFTP config create dto fields to entity
//
//Params
//	dto - TFTP config create dto
//	entity - dest TFTP config entity
func MapTFTPServerCreateDtoToEntity(dto dtos.TFTPServerCreateDto, entity *domain.TFTPConfig) {
	entity.Port = dto.Port
	entity.Address = dto.Address
	entity.Enabled = dto.Enabled
}

//MapTFTPServerUpdateDtoToEntity writes TFTP config update dto fields to entity
//
//Params
//	dto - TFTP config update dto
//	entity - dest TFTP config entity
func MapTFTPServerUpdateDtoToEntity(dto dtos.TFTPServerUpdateDto, entity *domain.TFTPConfig) {
	entity.Port = dto.Port
	entity.Address = dto.Address
	entity.Enabled = dto.Enabled
}
