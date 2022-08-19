package services

import (
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//HTTPLogService service structure for HTTPLog entity
type HTTPLogService struct {
	*GenericService[
		dtos.HTTPLogDto,
		dtos.HTTPLogDto,
		dtos.HTTPLogDto,
		domain.HTTPLog]
}

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
	genericSerice, err := NewGenericService[dtos.HTTPLogDto, dtos.HTTPLogDto, dtos.HTTPLogDto](rep, log)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error constructing http log service")
	}
	return HTTPLogService{
		genericSerice,
	}, nil
}
