package mappers

import (
	"fmt"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"
)

//MapDtoToEntity map a DTO to its corresponding entity
//Params
//	dto - DTO struct
//	entity - pointer to the entity
//Return
//  error - if error occurs return error, otherwise nil
func MapDtoToEntity(dto interface{}, entity interface{}) error {
	switch dto.(type) {
	// EthernetSwitch
	case dtos.EthernetSwitchCreateDto:
		MapEthernetSwitchCreateDto(dto.(dtos.EthernetSwitchCreateDto), entity.(*domain.EthernetSwitch))
	case dtos.EthernetSwitchUpdateDto:
		MapEthernetSwitchUpdateDto(dto.(dtos.EthernetSwitchUpdateDto), entity.(*domain.EthernetSwitch))
	//EthernetSwitchPort
	case dtos.EthernetSwitchPortCreateDto:
		MapEthernetSwitchPortCreateDto(dto.(dtos.EthernetSwitchPortCreateDto), entity.(*domain.EthernetSwitchPort))
	case dtos.EthernetSwitchPortUpdateDto:
		MapEthernetSwitchPortUpdateDto(dto.(dtos.EthernetSwitchPortUpdateDto), entity.(*domain.EthernetSwitchPort))
	default:
		return fmt.Errorf("[mapper]: Can't find route for map dto %+v to entity %+v", dto, entity)
	}
	return nil
}

//MapEntityToDto map entity to DTO
//Params
//	entity - Entity struct
//	dto - dest DTO
//Return
//  error - if error occurs return error, otherwise nil
func MapEntityToDto(entity interfaces.IEntityModel, dto interface{}) error {
	switch entity.(type) {
	// EthernetSwitch
	case domain.EthernetSwitch:
		MapEthernetSwitchToDto(entity.(domain.EthernetSwitch), dto.(*dtos.EthernetSwitchDto))
	//EthernetSwitchPort
	case domain.EthernetSwitchPort:
		MapEthernetSwitchPortToDto(entity.(domain.EthernetSwitchPort), dto.(*dtos.EthernetSwitchPortDto))
	//HttpLog
	case domain.HTTPLog:
		MapHTTPLogEntityToDto(entity.(domain.HTTPLog), dto.(*dtos.HTTPLogDto))
	//AppLog
	case domain.AppLog:
		MapAppLogEntityToDto(entity.(domain.AppLog), dto.(*dtos.AppLogDto))
	default:
		return fmt.Errorf("[mapper]: Can't find route for map entity %+v to dto %+v", entity, dto)
	}
	return nil
}
