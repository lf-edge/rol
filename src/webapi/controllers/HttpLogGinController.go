package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/webapi"
)

//HTTPLogGinController HTTP log controller structure for domain.HTTPLog
type HTTPLogGinController struct {
	service *services.HTTPLogService
	logger  *logrus.Logger
}

//RegisterHTTPLogController registers controller for getting http logs via api on path /api/v1/httplog/
func RegisterHTTPLogController(controller *HTTPLogGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/log/http/", controller.GetList)
	groupRoute.GET("/log/http/:id", controller.GetByID)
}

//NewHTTPLogGinController HTTP log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for http logs
func NewHTTPLogGinController(service *services.HTTPLogService, log *logrus.Logger) *HTTPLogGinController {
	logContr := &HTTPLogGinController{
		service: service,
		logger:  log,
	}
	return logContr
}

//GetList get list of http logs with search and pagination
//	Params
//	ctx - gin context
// @Summary Gets paginated list of http logs
// @version	1.0
// @Tags	log
// @Accept	json
// @Produce	json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"searchable value in entity"
// @param	page			query	int		false	"page number"
// @param	pageSize		query	int		false	"number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.HTTPLogDto]
// @Failure	500		"Internal Server Error"
// @router /log/http/ [get]
func (h *HTTPLogGinController) GetList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "CreatedAt", "desc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := h.service.GetList(ctx, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetByID get http log by id
//	Params
//	ctx - gin context
// @Summary	Gets http log by id
// @version	1.0
// @Tags	log
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"log id"
// @Success	200		{object}	dtos.HTTPLogDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /log/http/{id} [get]
func (h *HTTPLogGinController) GetByID(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := h.service.GetByID(ctx, id)
	handleWithData(ctx, err, dto)
}
