package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//HTTPLogService service structure for HTTPLog entity
type HTTPLogService struct {
	repository interfaces.IGenericRepository[uuid.UUID, domain.HTTPLog]
	logger     *logrus.Logger
}

//NewHTTPLogService preparing domain.HTTPLog repository for passing it in DI
//Params
//	repo - generic repository with HTTP log instantiated
//	log - logrus logger
//Return
//	New http log service
func NewHTTPLogService(repo interfaces.IGenericRepository[uuid.UUID, domain.HTTPLog], log *logrus.Logger) (*HTTPLogService, error) {
	return &HTTPLogService{
		repository: repo,
		logger:     log,
	}, nil
}

//GetList Get list of HTTP logs dto elements with search in all fields and pagination
//
//Params:
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return:
//	dtos.PaginatedItemsDto[dtos.HTTPLogDto] - paginated items
//	error - if an error occurs, otherwise nil
func (h *HTTPLogService) GetList(ctx context.Context, search string, orderBy string, orderDirection string, page int, pageSize int) (dtos.PaginatedItemsDto[dtos.HTTPLogDto], error) {
	return GetList[dtos.HTTPLogDto](ctx, h.repository, search, orderBy, orderDirection, page, pageSize)
}

//GetByID Get HTTP log dto by ID
//
//Params:
//	ctx - context
//	id - entity id
//Return:
//	dtos.HTTPLogDto - http log dto
//	error - if an error occurs, otherwise nil
func (h *HTTPLogService) GetByID(ctx context.Context, id uuid.UUID) (dtos.HTTPLogDto, error) {
	return GetByID[dtos.HTTPLogDto](ctx, h.repository, id, nil)
}
