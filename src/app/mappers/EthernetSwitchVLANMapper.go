package mappers

import (
	"github.com/google/uuid"
	"rol/domain"
	"rol/dtos"
	"strings"
)

func uuidSliceToString(slice []uuid.UUID) (out string) {
	for _, id := range slice {
		out += id.String() + ";"
	}
	return out[:len(out)-1]
}

func uuidsStringToSlice(IDs string) (out []uuid.UUID) {
	stringSlice := strings.Split(IDs, ";")
	for _, stringUUID := range stringSlice {
		id, _ := uuid.Parse(stringUUID)
		out = append(out, id)
	}
	return out
}

//MapEthernetSwitchVLANCreateDto writes ethernet switch port create dto fields to entity
//Params
//	dto - ethernet switch port create dto
//	entity - dest ethernet switch port entity
func MapEthernetSwitchVLANCreateDto(dto dtos.EthernetSwitchVLANCreateDto, entity *domain.EthernetSwitchVLAN) {
	entity.VlanID = dto.VlanID
	if dto.TaggedPorts != nil && len(dto.TaggedPorts) > 0 {
		entity.TaggedPorts = uuidSliceToString(dto.TaggedPorts)
	}
	if dto.UntaggedPorts != nil && len(dto.UntaggedPorts) > 0 {
		entity.UntaggedPorts = uuidSliceToString(dto.UntaggedPorts)
	}
}

//MapEthernetSwitchVLANUpdateDto writes ethernet switch port update dto fields to entity
//Params
//	dto - ethernet switch port update dto
//	entity - dest ethernet switch port entity
func MapEthernetSwitchVLANUpdateDto(dto dtos.EthernetSwitchVLANUpdateDto, entity *domain.EthernetSwitchVLAN) {
	if dto.TaggedPorts != nil && len(dto.TaggedPorts) > 0 {
		entity.TaggedPorts = uuidSliceToString(dto.TaggedPorts)
	} else {
		entity.TaggedPorts = ""
	}
	if dto.UntaggedPorts != nil && len(dto.UntaggedPorts) > 0 {
		entity.UntaggedPorts = uuidSliceToString(dto.UntaggedPorts)
	} else {
		entity.UntaggedPorts = ""
	}
}

//MapEthernetSwitchVLANToDto writes ethernet switch port entity to dto
//Params
//	entity - ethernet switch port entity
//	dto - dest ethernet switch port dto
func MapEthernetSwitchVLANToDto(entity domain.EthernetSwitchVLAN, dto *dtos.EthernetSwitchVLANDto) {
	dto.ID = entity.ID
	dto.CreatedAt = entity.CreatedAt
	dto.UpdatedAt = entity.UpdatedAt
	dto.VlanID = entity.VlanID
	dto.EthernetSwitchID = entity.EthernetSwitchID
	if entity.TaggedPorts != "" {
		dto.TaggedPorts = uuidsStringToSlice(entity.TaggedPorts)
	} else {
		dto.TaggedPorts = []uuid.UUID{}
	}
	if entity.UntaggedPorts != "" {
		dto.UntaggedPorts = uuidsStringToSlice(entity.UntaggedPorts)
	} else {
		dto.UntaggedPorts = []uuid.UUID{}
	}
}
