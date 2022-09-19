package mappers

import (
	"net"
	"rol/domain"
	"rol/dtos"
)

//MapHostNetworkVlanToDto map HostNetworkVlan entity to dto
func MapHostNetworkVlanToDto(entity domain.HostNetworkVlan, dto *dtos.HostNetworkVlanDto) {
	dto.Name = entity.Name
	dto.VlanID = entity.VlanID
	for _, addr := range entity.Addresses {
		dto.Addresses = append(dto.Addresses, addr.String())
	}
	dto.Parent = entity.Parent
}

//MapHostNetworkCreateDtoToEntity map HostNetworkCreateDto dto to entity
func MapHostNetworkCreateDtoToEntity(dto dtos.HostNetworkVlanCreateDto, entity *domain.HostNetworkVlan) {
	entity.VlanID = dto.VlanID
	entity.Parent = dto.Parent
	entity.Type = "vlan"
	for _, addr := range dto.Addresses {
		ip, address, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		address.IP = ip
		entity.Addresses = append(entity.Addresses, *address)
	}
}

//MapHostNetworkUpdateDtoToEntity map HostNetworkUpdateDto dto to entity
func MapHostNetworkUpdateDtoToEntity(dto dtos.HostNetworkVlanUpdateDto, entity *domain.HostNetworkVlan) {
	entity.Addresses = []net.IPNet{}
	for _, addr := range dto.Addresses {
		ip, address, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		address.IP = ip
		entity.Addresses = append(entity.Addresses, *address)
	}
}
