package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/webapi"
)

//AppLogGinController Application log GIN controller constructor
type AppLogGinController struct {
	service *services.AppLogService
	logger  *logrus.Logger
}

//RegisterAppLogController registers controller for getting application logs via api on path /api/v1/applog/
func RegisterAppLogController(controller *AppLogGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/log/app/", controller.GetList)
	groupRoute.GET("/log/app/:id", controller.GetByID)
}

//NewAppLogGinController app log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for app logs
func NewAppLogGinController(service *services.AppLogService, log *logrus.Logger) *AppLogGinController {
	logContr := &AppLogGinController{
		service: service,
		logger:  log,
	}
	return logContr
}

//GetList get list of logs with search and pagination
//	Params
//	ctx - gin context
// @Summary Gets paginated list of app logs
// @version	1.0
// @Tags	log
// @Accept	json
// @Produce	json
// @param	orderBy			query		string	false	"Order by field"
// @param	orderDirection	query		string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query		string	false	"searchable value in entity"
// @param	page			query		int		false	"page number"
// @param	pageSize		query		int		false	"number of entities per page"
// @Success	200				{object}	dtos.PaginatedItemsDto[dtos.AppLogDto]
// @Failure	500				"Internal Server Error"
// @router /log/app/ [get]
func (a *AppLogGinController) GetList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "CreatedAt", "desc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := a.service.GetList(ctx, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetByID get log by id
//	Params
//	ctx - gin context
// @Summary Gets http app by id
// @version	1.0
// @Tags	log
// @Accept	json
// @Produce	json
// @param	id		path					string			true	"log id"
// @Success	200		{object}				dtos.AppLogDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /log/app/{id} [get]
func (a *AppLogGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := a.service.GetByID(ctx, id)
	handleWithData(ctx, err, dto)
}
