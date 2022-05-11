package services

import (
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//NewHTTPLogService preparing domain.HTTPLog repository for passing it in DI
//Params
//	rep - generic repository with http log instantiated
//	log - logrus logger
//Return
//	New http log service
func NewHTTPLogService(rep interfaces.IGenericRepository[domain.HTTPLog], log *logrus.Logger) (interfaces.IGenericService[
	dtos.HTTPLogDto,
	dtos.HTTPLogDto,
	dtos.HTTPLogDto,
	domain.HTTPLog], error) {
	return NewGenericService[dtos.HTTPLogDto, dtos.HTTPLogDto, dtos.HTTPLogDto](rep, log)
}
