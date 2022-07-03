package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/webapi"

	"github.com/sirupsen/logrus"
)

//EthernetSwitchGinController ethernet switch GIN controller constructor
type EthernetSwitchGinController struct {
	GinGenericController[dtos.EthernetSwitchDto,
		dtos.EthernetSwitchCreateDto,
		dtos.EthernetSwitchUpdateDto,
		domain.EthernetSwitch]
}

//RegisterSwitchController registers controller for the switches on path /api/v1/switch/
func RegisterSwitchController(controller *EthernetSwitchGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/switch/", controller.GetList)
	groupRoute.GET("/switch/:id", controller.GetByID)
	groupRoute.POST("/switch", controller.Create)
	groupRoute.PUT("/switch/:id", controller.Update)
	groupRoute.DELETE("/switch/:id", controller.Delete)
	groupRoute.GET("/switch/models", controller.GetSupportedModels)
}

//GetList get list of switches with search and pagination
//	Params
//	ctx - gin context
// @Summary Gets paginated list of switches
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce json
// @param	 orderBy			path	string	false	"Order by field"
// @param	 orderDirection	path	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	 search			 path	string	false	"searchable value in entity"
// @param	 page			 path	int		false	"page number"
// @param	 pageSize		 path	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.EthernetSwitchDto}} "desc"
// @router /switch/ [get]
func (e *EthernetSwitchGinController) GetList(ctx *gin.Context) {
	e.GinGenericController.GetList(ctx)
}

//GetByID get switch by id
//	Params
//	ctx - gin context
// @Summary Gets ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @param	 id	path	string		true	"ethernet switch id"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.EthernetSwitchDto}
// @router /switch/{id} [get]
func (e *EthernetSwitchGinController) GetByID(ctx *gin.Context) {
	e.GinGenericController.GetByID(ctx)
}

//Create new switch
//	Params
//	ctx - gin context
// @Summary Creates new ethernet switch
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @Param request body dtos.EthernetSwitchCreateDto true "ethernet switch fields"
// @Success 200 {object} dtos.ResponseDataDto
// @router /switch/ [post]
func (e *EthernetSwitchGinController) Create(ctx *gin.Context) {
	e.GinGenericController.Create(ctx)
}

//Update switch in database by id
//	Params
//	ctx - gin context
// @Summary Updates ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @Param request body dtos.EthernetSwitchUpdateDto true "ethernet switch fields"
// @Success 200 {object} dtos.ResponseDto
// @router /switch/{id} [put]
func (e *EthernetSwitchGinController) Update(ctx *gin.Context) {
	e.GinGenericController.Update(ctx)
}

//Delete soft deleting switch in database
//	Params
//	ctx - gin context
// @Summary Delete ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @param	 id	path	string		true	"ethernet switch id"
// @Success 200 {object} dtos.ResponseDto
// @router /switch/{id} [delete]
func (e *EthernetSwitchGinController) Delete(ctx *gin.Context) {
	e.GinGenericController.Delete(ctx)
}

//GetSupportedModels Get supported switch models
//	Params
//	ctx - gin context
// @Summary Get ethernet switch supported models
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @Success 200 {object} []dtos.EthernetSwitchModelDto
// @router /switch/models [get]
func (e *EthernetSwitchGinController) GetSupportedModels(ctx *gin.Context) {
	service := e.GinGenericController.service.(*services.EthernetSwitchService)
	modelsDtoSlice := service.GetSupportedModels()
	ctx.JSON(http.StatusOK, modelsDtoSlice)
}

//NewEthernetSwitchGinController ethernet switch controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for ethernet switches
func NewEthernetSwitchGinController(service interfaces.IGenericService[dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], log *logrus.Logger) *EthernetSwitchGinController {
	genContr := NewGinGenericController(service, log)
	switchContr := &EthernetSwitchGinController{
		GinGenericController: *genContr,
	}
	return switchContr
}
