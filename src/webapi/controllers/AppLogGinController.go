package controllers

import (
	"github.com/gin-gonic/gin"
	"rol/app/interfaces"
	"rol/domain"
	"rol/dtos"
	"rol/webapi"

	"github.com/sirupsen/logrus"
)

//AppLogGinController Application log GIN controller constructor
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
	groupRoute.GET("/log/app/:id", controller.GetByID)
}

//GetList get list of logs with search and pagination
//	Params
//	ctx - gin context
// @Summary Gets paginated list of app logs
// @version 1.0
// @Tags app log
// @Accept  json
// @Produce  json
// @param	 orderBy		 query	string	false	"Order by field"
// @param	 orderDirection		 query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	 search			 query	string	false	"searchable value in entity"
// @param	 page			 query	int		false	"page number"
// @param	 pageSize		 query	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.AppLogDto}} ""
// @router /log/app/ [get]
func (a *AppLogGinController) GetList(ctx *gin.Context) {
	a.GinGenericController.GetList(ctx)
}

//GetByID get log by id
//	Params
//	ctx - gin context
// @Summary Gets http app by id
// @version 1.0
// @Tags app log
// @Accept  json
// @Produce  json
// @param	 id	path	string		true	"log id"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.AppLogDto}
// @router /log/app/{id} [get]
func (a *AppLogGinController) GetByID(ctx *gin.Context) {
	a.GinGenericController.GetByID(ctx)
}

//NewAppLogGinController app log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for app logs
func NewAppLogGinController(service interfaces.IGenericService[dtos.AppLogDto,
	dtos.AppLogDto,
	dtos.AppLogDto,
	domain.AppLog], log *logrus.Logger) *AppLogGinController {
	genContr := NewGinGenericController(service, log)
	logContr := &AppLogGinController{
		GinGenericController: *genContr,
	}
	return logContr
}
