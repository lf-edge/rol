package mappers

import (
	"rol/domain"
	"rol/dtos"
)

//MapDHCP4ServerToDto writes dhcp v4 config fields to dto
//
//Params:
//	entity - DHCP v4 config entity
//	*dto - DHCP v4 server dto
func MapDHCP4ServerToDto(entity domain.DHCP4Config, dto *dtos.DHCP4ServerDto) {
	dto.ID = entity.ID
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
	dto.Interface = entity.Interface
	dto.DNS = entity.DNS
	dto.NTP = entity.NTP
	dto.Range = entity.Range
	dto.Mask = entity.Mask
	dto.ServerID = entity.ServerID
	dto.Gateway = entity.Gateway
	dto.Port = entity.Port
	dto.Enabled = entity.Enabled
	dto.LeaseTime = entity.LeaseTime
}

//MapDHCP4ServerCreateDtoToEntity writes dhcp v4 create dto fields to entity
//
//Params:
// 	dto - DHCP v4 server create dto
//	entity - DHCP v4 config entity
func MapDHCP4ServerCreateDtoToEntity(dto dtos.DHCP4ServerCreateDto, entity *domain.DHCP4Config) {
	entity.DNS = dto.DNS
	entity.Interface = dto.Interface
	entity.NTP = dto.NTP
	entity.Range = dto.Range
	entity.Mask = dto.Mask
	entity.ServerID = dto.ServerID
	entity.Gateway = dto.Gateway
	entity.Port = dto.Port
	entity.Enabled = dto.Enabled
	entity.LeaseTime = dto.LeaseTime
}

//MapDHCP4ServerUpdateDtoToEntity writes dhcp v4 update dto fields to entity
//
//Params:
// 	dto - DHCP v4 server create dto
//	entity - DHCP v4 config entity
func MapDHCP4ServerUpdateDtoToEntity(dto dtos.DHCP4ServerUpdateDto, entity *domain.DHCP4Config) {
	entity.DNS = dto.DNS
	entity.NTP = dto.NTP
	entity.Port = dto.Port
	entity.Enabled = dto.Enabled
	entity.LeaseTime = dto.LeaseTime
}

//MapDHCP4LeaseToDto writes dhcp v4 lease fields to dto
//
//Params:
//	entity - DHCP v4 lease entity
//	*dto - DHCP v4 lease dto
func MapDHCP4LeaseToDto(entity domain.DHCP4Lease, dto *dtos.DHCP4LeaseDto) {
	dto.IP = entity.IP
	dto.MAC = entity.MAC
	dto.Expires = entity.Expires
	dto.UpdatedAt = entity.UpdatedAt
	dto.CreatedAt = entity.CreatedAt
	dto.ID = entity.ID
}

//MapDHCP4LeaseCreateDtoToEntity writes dhcp v4 lease create dto fields to lease entity
//
//Params:
// 	dto - DHCP v4 server lease create dto
//	entity - DHCP v4 lease entity
func MapDHCP4LeaseCreateDtoToEntity(dto dtos.DHCP4LeaseCreateDto, entity *domain.DHCP4Lease) {
	entity.IP = dto.IP
	entity.MAC = dto.MAC
	entity.Expires = dto.Expires
}

//MapDHCP4LeaseUpdateDtoToEntity writes dhcp v4 lease update dto fields to lease entity
//
//Params:
// 	dto - DHCP v4 server lease create dto
//	entity - DHCP v4 lease entity
func MapDHCP4LeaseUpdateDtoToEntity(dto dtos.DHCP4LeaseUpdateDto, entity *domain.DHCP4Lease) {
	entity.IP = dto.IP
	entity.MAC = dto.MAC
	entity.Expires = dto.Expires
}
