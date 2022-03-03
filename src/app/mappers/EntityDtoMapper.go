package mappers

import (
	"rol/app/interfaces"
	"rol/domain/entities"
	"rol/dtos"
)

func GetEmptyEntity(dto interfaces.IEntityDtoModel) interfaces.IEntityModel {
	switch dto.(type) {
	case *dtos.EthernetSwitchCreateDto, *dtos.EthernetSwitchUpdateDto, *dtos.EthernetSwitchDto:
		return &entities.EthernetSwitch{}
	case *dtos.EthernetSwitchPortCreateDto, *dtos.EthernetSwitchPortUpdateDto, *dtos.EthernetSwitchPortDto:
		return &entities.EthernetSwitchPort{}
	default:
		return nil
	}
}

func GetEntityEmptyArray(dto interface{}) interface{} {
	switch dto.(type) {
	case *[]*dtos.EthernetSwitchDto:
		return &[]*entities.EthernetSwitch{}
	case *[]*dtos.EthernetSwitchPortDto:
		return &[]*entities.EthernetSwitchPort{}
	}
	return nil
}

func Map(source interface{}, dest interface{}) {
	switch source.(type) {
	// EthernetSwitch
	case *dtos.EthernetSwitchCreateDto:
		MapEthernetSwitchCreateDto(source.(*dtos.EthernetSwitchCreateDto), dest.(*entities.EthernetSwitch))
	case *dtos.EthernetSwitchUpdateDto:
		MapEthernetSwitchUpdateDto(source.(*dtos.EthernetSwitchUpdateDto), dest.(*entities.EthernetSwitch))
	case *entities.EthernetSwitch:
		MapEthernetSwitchDto(source.(*entities.EthernetSwitch), dest.(*dtos.EthernetSwitchDto))
	case *[]*entities.EthernetSwitch:
		MapEthernetSwitchArrayDto(source.(*[]*entities.EthernetSwitch), dest.(*[]*dtos.EthernetSwitchDto))
	//EthernetSwitchPort
	case *dtos.EthernetSwitchPortCreateDto:
		MapEthernetSwitchPortCreateDto(source.(*dtos.EthernetSwitchPortCreateDto), dest.(*entities.EthernetSwitchPort))
	case *dtos.EthernetSwitchPortUpdateDto:
		MapEthernetSwitchPortUpdateDto(source.(*dtos.EthernetSwitchPortUpdateDto), dest.(*entities.EthernetSwitchPort))
	case *entities.EthernetSwitchPort:
		MapEthernetSwitchPortDto(source.(*entities.EthernetSwitchPort), dest.(*dtos.EthernetSwitchPortDto))
	case *[]*entities.EthernetSwitchPort:
		MapEthernetSwitchPortArrayDto(source.(*[]*entities.EthernetSwitchPort), dest.(*[]*dtos.EthernetSwitchPortDto))
	default:
		return
	}
}
