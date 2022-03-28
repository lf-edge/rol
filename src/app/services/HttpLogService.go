package services

import (
	"rol/app/interfaces/generic"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//NewHttpLogService preparing domain.HttpLog repository for passing it in DI
//Params
//	rep - generic repository with http log instantiated
//	log - logrus logger
//Return
//	New http log service
func NewHttpLogService(rep generic.IGenericRepository[domain.HttpLog], log *logrus.Logger) (generic.IGenericService[
	dtos.HttpLogDto,
	dtos.HttpLogDto,
	dtos.HttpLogDto,
	domain.HttpLog], error) {
	return NewGenericService[dtos.HttpLogDto, dtos.HttpLogDto, dtos.HttpLogDto](rep, log)
}
