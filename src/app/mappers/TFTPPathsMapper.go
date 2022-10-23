// Package mappers uses for entity <--> dto conversions
package mappers

import (
	"github.com/google/uuid"
	"rol/domain"
	"rol/dtos"
)

//MapTFTPPathRatioToDto writes TFTP path ratio entity to dto
//Params
//	entity - TFTP path entity
//	dto - dest TFTP path dto
func MapTFTPPathRatioToDto(entity domain.TFTPPathRatio, dto *dtos.TFTPPathDto) {
	mapEntityToBaseDto[uuid.UUID](entity, &dto.BaseDto)
	dto.ActualPath = entity.ActualPath
	dto.VirtualPath = entity.VirtualPath
}

//MapTFTPPathCreateDtoToEntity writes TFTP path create dto fields to entity
//Params
//	dto -TFTP path create dto
//	entity - dest TFTP path entity
func MapTFTPPathCreateDtoToEntity(dto dtos.TFTPPathCreateDto, entity *domain.TFTPPathRatio) {
	entity.ActualPath = dto.ActualPath
	entity.VirtualPath = dto.VirtualPath
}

//MapTFTPPathUpdateDtoToEntity writes TFTP path update dto fields to entity
//Params
//	dto -TFTP path create dto
//	entity - dest TFTP path entity
func MapTFTPPathUpdateDtoToEntity(dto dtos.TFTPPathUpdateDto, entity *domain.TFTPPathRatio) {
	entity.ActualPath = dto.ActualPath
	entity.VirtualPath = dto.VirtualPath
}
