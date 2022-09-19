package mappers

import (
	"net"
	"rol/domain"
	"rol/dtos"
)

//MapHostNetworkBridgeToDto map HostNetworkVlan entity to dto
func MapHostNetworkBridgeToDto(entity domain.HostNetworkBridge, dto *dtos.HostNetworkBridgeDto) {
	dto.Name = entity.Name
	for _, addr := range entity.Addresses {
		dto.Addresses = append(dto.Addresses, addr.String())
	}
	dto.Slaves = entity.Slaves
}

//MapHostNetworkBridgeCreateDtoToEntity map HostNetworkCreateDto dto to entity
func MapHostNetworkBridgeCreateDtoToEntity(dto dtos.HostNetworkBridgeCreateDto, entity *domain.HostNetworkBridge) {
	entity.Type = "bridge"
	for _, addr := range dto.Addresses {
		ip, address, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		address.IP = ip
		entity.Addresses = append(entity.Addresses, *address)
	}
	entity.Slaves = dto.Slaves
	entity.Name = dto.Name
}

//MapHostNetworkBridgeUpdateDtoToEntity map HostNetworkUpdateDto dto to entity
func MapHostNetworkBridgeUpdateDtoToEntity(dto dtos.HostNetworkBridgeUpdateDto, entity *domain.HostNetworkBridge) {
	entity.Addresses = []net.IPNet{}
	for _, addr := range dto.Addresses {
		ip, address, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		address.IP = ip
		entity.Addresses = append(entity.Addresses, *address)
	}
	entity.Slaves = dto.Slaves
}
