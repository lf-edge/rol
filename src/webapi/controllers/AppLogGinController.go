package controllers

import (
	"github.com/gin-gonic/gin"
	"rol/app/interfaces/generic"
	"rol/domain"
	"rol/dtos"
	"rol/webapi"

	"github.com/sirupsen/logrus"
)

//NewAppLogGinController Application log GIN controller constructor
type AppLogGinController struct {
	GinGenericController[dtos.AppLogDto,
		dtos.AppLogDto,
		dtos.AppLogDto,
		domain.AppLog]
}

//RegisterAppLogController registers controller for getting application logs via api on path /api/v1/applog/
func RegisterAppLogController(controller *AppLogGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/log/app/", controller.GetList)
	groupRoute.GET("/log/app/:id", controller.GetById)
}

// @Summary Gets paginated list of app logs
// @version 1.0
// @Tags app log
// @Accept  json
// @Produce  json
// @param 	orderBy			path	string	false	"Order by field"
// @param 	orderDirection	path	string	false	"'asc' or 'desc' for ascending or descending order"
// @param 	search 			path	string	false	"searchable value in entity"
// @param 	page 			path	int		false	"page number"
// @param 	pageSize 		path	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.AppLogDto}} ""
// @router /log/app/ [get]
func (c *AppLogGinController) GetList(ctx *gin.Context) {
	c.GinGenericController.GetList(ctx)
}

// @Summary Gets http app by id
// @version 1.0
// @Tags app log
// @Accept  json
// @Produce  json
// @param 	id	path	string		true	"log id"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.AppLogDto}
// @router /log/app/{id} [get]
func (c *AppLogGinController) GetById(ctx *gin.Context) {
	c.GinGenericController.GetById(ctx)
}

//NewEthernetSwitchGinController ppp log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for app logs
func NewAppLogGinController(service generic.IGenericService[dtos.AppLogDto,
	dtos.AppLogDto,
	dtos.AppLogDto,
	domain.AppLog], log *logrus.Logger) *AppLogGinController {
	genContr := NewGinGenericController(service, log)
	logContr := &AppLogGinController{
		GinGenericController: *genContr,
	}
	return logContr
}
