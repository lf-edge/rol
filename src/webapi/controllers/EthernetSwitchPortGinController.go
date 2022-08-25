package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
	"strconv"
)

//EthernetSwitchPortGinController Ethernet switch port API controller for domain.EthernetSwitchPort entity
type EthernetSwitchPortGinController struct {
	service *services.EthernetSwitchService
	logger  *logrus.Logger
}

//NewEthernetSwitchPortGinController HTTP log controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for http logs
func NewEthernetSwitchPortGinController(service *services.EthernetSwitchService, log *logrus.Logger) *EthernetSwitchPortGinController {
	ethernetSwitchPortController := &EthernetSwitchPortGinController{
		service: service,
		logger:  log,
	}
	return ethernetSwitchPortController
}

//RegisterEthernetSwitchPortController registers controller for getting ethernet switch ports via api
func RegisterEthernetSwitchPortController(controller *EthernetSwitchPortGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/ethernet-switch/:id/port/", controller.GetPorts)
	groupRoute.GET("/ethernet-switch/:id/port/:portID", controller.GetPortByID)
	groupRoute.POST("/ethernet-switch/:id/port/", controller.CreatePort)
	groupRoute.PUT("/ethernet-switch/:id/port/:portID", controller.UpdatePort)
	groupRoute.DELETE("/ethernet-switch/:id/port/:portID", controller.DeletePort)
}

//GetPortByID Get ethernet switch port by id
//Params
//	ctx - gin context
// @Summary Gets ethernet switch port by id
// @version	1.0
// @Tags	ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Ethernet switch ID"
// @param	portID	path		string		true	"Ethernet switch port ID"
// @Success	200		{object}	dtos.EthernetSwitchPortDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/port/{portID} [get]
func (e *EthernetSwitchPortGinController) GetPortByID(ctx *gin.Context) {
	strID := ctx.Param("portID")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.GetPortByID(ctx, switchID, id)
	handleWithData(ctx, err, dto)
}

//GetPorts Get list of elements with search and pagination
//Params
//	ctx - gin context
// @Summary	Gets paginated list of ethernet switch ports
// @version	1.0
// @Tags	ethernet-switch
// @Accept  json
// @Produce	json
// @param	id				path	string	true	"Ethernet switch ID"
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.EthernetSwitchPortDto]
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/port/ [get]
func (e *EthernetSwitchPortGinController) GetPorts(ctx *gin.Context) {
	orderBy := ctx.DefaultQuery("orderBy", "Name")
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
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	paginatedList, err := e.service.GetPorts(ctx, switchID, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	handleWithData(ctx, err, paginatedList)
}

//CreatePort new ethernet switch port
//	Params
//	ctx - gin context
// @Summary Creates new ethernet switch port
// @version 1.0
// @Tags ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path		string		                        true	"Ethernet switch ID"
// @Param	request	body		dtos.EthernetSwitchPortCreateDto    true    "Ethernet switch port fields"
// @Success	200		{object}	dtos.EthernetSwitchPortDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/port/ [post]
func (e *EthernetSwitchPortGinController) CreatePort(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchPortCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.CreatePort(ctx, switchID, reqDto)
	handleWithData(ctx, err, dto)
}

//UpdatePort Update Ethernet switch port by id
//	Params
//	ctx - gin context
// @Summary Updates ethernet switch port by id
// @version     1.0
// @Tags        ethernet-switch
// @Accept      json
// @Produce     json
// @param       id          path    string		                        true	"Ethernet switch ID"
// @param       portID      path    string		                        true	"Ethernet switch port ID"
// @Param       request     body    dtos.EthernetSwitchPortUpdateDto    true    "Ethernet switch port fields"
// @Success		200		{object}	dtos.EthernetSwitchPortDto
// @Failure		400		{object}	dtos.ValidationErrorDto
// @Failure		404		"Not Found"
// @Failure		500		"Internal Server Error"
// @router /ethernet-switch/{id}/port/{portID} [put]
func (e *EthernetSwitchPortGinController) UpdatePort(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchPortUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	portID, err := parseUUIDParam(ctx, "portID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.UpdatePort(ctx, switchID, portID, reqDto)
	handleWithData(ctx, err, dto)
}

//DeletePort deleting ethernet switch port
//	Params
//	ctx - gin context
// @Summary Delete ethernet switch port by id
// @version 1.0
// @Tags ethernet-switch
// @Accept	json
// @Produce	json
// @param	id		path					string	true	"Ethernet switch ID"
// @param	portID	path					string	true	"Ethernet switch port ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/port/{portID} [delete]
func (e *EthernetSwitchPortGinController) DeletePort(ctx *gin.Context) {
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	portID, err := parseUUIDParam(ctx, "portID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = e.service.DeletePort(ctx, switchID, portID)
	handle(ctx, err)
}
