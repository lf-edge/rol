package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapProjectToDto map Project entity to dto
func MapProjectToDto(entity domain.Project, dto *dtos.ProjectDto) {
	dto.ID = entity.ID
	dto.Name = entity.Name
	dto.Subnet = entity.Subnet
	dto.BridgeName = entity.BridgeName
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
	dto.DHCPServerID = entity.DHCPServerID
	dto.TFTPServerID = entity.TFTPServerID
}

//MapProjectCreateDtoToEntity map ProjectCreateDto dto to entity
func MapProjectCreateDtoToEntity(dto dtos.ProjectCreateDto, entity *domain.Project) {
	entity.Name = dto.Name
	entity.Subnet = dto.Subnet
}
