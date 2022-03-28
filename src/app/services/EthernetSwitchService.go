package services

import (
	"rol/app/interfaces/generic"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//NewEthernetSwitchService constructor for domain.EthernetSwitch service
//Params
//	rep - generic repository with domain.EthernetSwitch repository
//	log - logrus logger
//Return
//	New ethernet switch service
func NewEthernetSwitchService(rep generic.IGenericRepository[domain.EthernetSwitch], log *logrus.Logger) (generic.IGenericService[
	dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], error) {
	return NewGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto](rep, log)
}
