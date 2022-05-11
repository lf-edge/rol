package controllers

import (
	"github.com/gin-gonic/gin"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"
	"rol/webapi"

	"github.com/sirupsen/logrus"
)

type HTTPLogGinController struct {
	GinGenericController[dtos.HTTPLogDto,
		dtos.HTTPLogDto,
		dtos.HTTPLogDto,
		domain.HTTPLog]
}

//RegisterHTTPLogController registers controller for getting http logs via api on path /api/v1/httplog/
func RegisterHTTPLogController(controller *HTTPLogGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/log/http/", controller.GetList)
	groupRoute.GET("/log/http/:id", controller.GetById)
}

// @Summary Gets paginated list of http logs
// @version 1.0
// @Tags http log
// @Accept  json
// @Produce  json
// @param 	orderBy			path	string	false	"Order by field"
// @param 	orderDirection	path	string	false	"'asc' or 'desc' for ascending or descending order"
// @param 	search 			path	string	false	"searchable value in entity"
// @param 	page 			path	int		false	"page number"
// @param 	pageSize 		path	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.HTTPLogDto}} "desc"
// @router /log/http/ [get]
func (c *HTTPLogGinController) GetList(ctx *gin.Context) {
	c.GinGenericController.GetList(ctx)
}

// @Summary Gets http log by id
// @version 1.0
// @Tags http log
// @Accept  json
// @Produce  json
// @param 	id	path	string		true	"log id"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.HTTPLogDto}
// @router /log/http/{id} [get]
func (c *HTTPLogGinController) GetById(ctx *gin.Context) {
	c.GinGenericController.GetById(ctx)
}

//NewHTTPLogGinController HTTP log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for http logs
func NewHTTPLogGinController(service interfaces.IGenericService[dtos.HTTPLogDto,
	dtos.HTTPLogDto,
	dtos.HTTPLogDto,
	domain.HTTPLog], log *logrus.Logger) *HTTPLogGinController {
	genContr := NewGinGenericController(service, log)
	logContr := &HTTPLogGinController{
		GinGenericController: *genContr,
	}
	return logContr
}
