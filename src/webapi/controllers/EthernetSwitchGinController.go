package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
	"strconv"

	"github.com/sirupsen/logrus"
)

//EthernetSwitchGinController ethernet switch GIN controller constructor
type EthernetSwitchGinController struct {
	service *services.EthernetSwitchService
	logger  *logrus.Logger
}

//RegisterEthernetSwitchController registers controller for the switches on path /api/v1/ethernet-switch/
func RegisterEthernetSwitchController(controller *EthernetSwitchGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/ethernet-switch/", controller.GetList)
	groupRoute.GET("/ethernet-switch/models/", controller.GetSupportedModels)
	groupRoute.GET("/ethernet-switch/:id", controller.GetByID)
	groupRoute.POST("/ethernet-switch", controller.Create)
	groupRoute.PUT("/ethernet-switch/:id", controller.Update)
	groupRoute.DELETE("/ethernet-switch/:id", controller.Delete)
}

//NewEthernetSwitchGinController ethernet switch controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for ethernet switches
func NewEthernetSwitchGinController(service *services.EthernetSwitchService, log *logrus.Logger) *EthernetSwitchGinController {
	switchContr := &EthernetSwitchGinController{
		service: service,
		logger:  log,
	}
	return switchContr
}

//GetList get list of switches with search and pagination
//	Params
//	ctx - gin context
// @Summary Get paginated list of switches
// @version 1.0
// @Tags	ethernet-switch
// @Accept  json
// @Produce json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.EthernetSwitchDto]
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/ [get]
func (e *EthernetSwitchGinController) GetList(ctx *gin.Context) {
	orderBy := ctx.DefaultQuery("orderBy", "id")
	orderDirection := ctx.DefaultQuery("orderDirection", "asc")
	search := ctx.DefaultQuery("search", "")
	page := ctx.DefaultQuery("page", "1")
	pageInt64, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		pageInt64 = 1
	}
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt64, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		pageSizeInt64 = 10
	}
	paginatedList, err := e.service.GetList(ctx, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	handleWithData(ctx, err, paginatedList)
}

//GetByID get switch by id
//	Params
//	ctx - gin context
// @Summary	Get ethernet switch by id
// @version 1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Ethernet switch ID"
// @Success	200		{object}	dtos.EthernetSwitchDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id} [get]
func (e *EthernetSwitchGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.GetByID(ctx, id)
	handleWithData(ctx, err, dto)
}

//Create new switch
//	Params
//	ctx - gin context
// @Summary	Create new ethernet switch
// @version	1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.EthernetSwitchCreateDto	true	"Ethernet switch fields"
// @Success	200		{object}	dtos.EthernetSwitchDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/ [post]
func (e *EthernetSwitchGinController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := e.service.Create(ctx, reqDto)
	handleWithData(ctx, err, dto)
}

//Update switch by id
//	Params
//	ctx - gin context
// @Summary	Updates ethernet switch by id
// @version	1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Ethernet switch ID"
// @Param	request	body		dtos.EthernetSwitchUpdateDto true "Ethernet switch fields"
// @Success	200		{object}	dtos.EthernetSwitchDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id} [put]
func (e *EthernetSwitchGinController) Update(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := e.service.Update(ctx, reqDto, id)
	handleWithData(ctx, err, dto)
}

//Delete soft deleting switch in database
//	Params
//	ctx - gin context
// @Summary	Delete ethernet switch by id
// @version	1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path	string		true	"Ethernet switch ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id} [delete]
func (e *EthernetSwitchGinController) Delete(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	err = e.service.Delete(ctx, id)
	handle(ctx, err)
}

//GetSupportedModels Get supported switch models
//	Params
//	ctx - gin context
// @Summary	Get ethernet switch supported models
// @version	1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @Success	200		{object} []dtos.EthernetSwitchModelDto
// @router /ethernet-switch/models [get]
func (e *EthernetSwitchGinController) GetSupportedModels(ctx *gin.Context) {
	modelsDtoSlice := e.service.GetSupportedModels()
	ctx.JSON(http.StatusOK, modelsDtoSlice)
}
