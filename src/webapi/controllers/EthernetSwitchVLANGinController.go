package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//EthernetSwitchVLANGinController ethernet switch GIN controller constructor
type EthernetSwitchVLANGinController struct {
	service *services.EthernetSwitchService
	logger  *logrus.Logger
}

//NewEthernetSwitchVLANGinController ethernet switch VLAN controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for http logs
func NewEthernetSwitchVLANGinController(service *services.EthernetSwitchService, log *logrus.Logger) *EthernetSwitchVLANGinController {
	ethernetSwitchVLANController := &EthernetSwitchVLANGinController{
		service: service,
		logger:  log,
	}
	return ethernetSwitchVLANController
}

//RegisterEthernetSwitchVLANGinController registers controller for getting ethernet switch VLANs via api
func RegisterEthernetSwitchVLANGinController(controller *EthernetSwitchVLANGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/ethernet-switch/:id/vlan/", controller.GetList)
	groupRoute.GET("/ethernet-switch/:id/vlan/:vlanUUID", controller.GetByID)
	groupRoute.POST("/ethernet-switch/:id/vlan/", controller.Create)
	groupRoute.PUT("/ethernet-switch/:id/vlan/:vlanUUID", controller.Update)
	groupRoute.DELETE("/ethernet-switch/:id/vlan/", controller.Delete)
}

//GetList get list of switch VLANs with search and pagination
//	Params
//	ctx - gin context
// @Summary Get paginated list of switch VLANs
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce json
// @param 	 id 			 path   string  true "Ethernet switch ID"
// @param	 orderBy		 query	string	false	"Order by field"
// @param	 orderDirection	 query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	 search			 query	string	false	"Searchable value in entity"
// @param	 page			 query	int		false	"Page number"
// @param	 pageSize		 query	int		false	"Number of entities per page"
// @Success 200 {object} dtos.PaginatedItemsDto[dtos.EthernetSwitchVLANDto]
// @Failure		500		"Internal Server Error"
// @router /ethernet-switch/{id}/vlan [get]
func (e *EthernetSwitchVLANGinController) GetList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "VlanID", "asc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := e.service.GetVLANs(ctx, switchID, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetByID get switch VLAN by id
//	Params
//	ctx - gin context
// @Summary Get ethernet switch VLAN by id
// @version 1.0
// @Tags 	ethernet-switch
// @Accept  json
// @Produce json
// @param	id		path		string		true	"Ethernet switch ID"
// @param	vlanUUID	path		string		true	"Ethernet switch VLAN UUID"
// @Success 200 	{object} 	dtos.EthernetSwitchVLANDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/vlan/{vlanUUID} [get]
func (e *EthernetSwitchVLANGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("vlanUUID")
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
	dto, err := e.service.GetVLANByID(ctx, switchID, id)
	handleWithData(ctx, err, dto)
}

//Create new switch VLAN
//	Params
//	ctx - gin context
// @Summary Create new ethernet switch VLAN
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce json
// @Param 	id 		path 		string true "Ethernet switch ID"
// @Param 	request body 		dtos.EthernetSwitchVLANCreateDto true "Ethernet switch VLAN fields"
// @Success 200 	{object} 	dtos.EthernetSwitchVLANDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/vlan [post]
func (e *EthernetSwitchVLANGinController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchVLANCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.CreateVLAN(ctx, switchID, reqDto)
	handleWithData(ctx, err, dto)
}

//Update switch VLAN by id
//	Params
//	ctx - gin context
// @Summary Updates ethernet switch VLAN by id
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce  json
// @param id path string true "Ethernet switch ID"
// @param vlanUUID path string true "Ethernet switch VLAN UUID"
// @Param request body dtos.EthernetSwitchVLANUpdateDto true "Ethernet switch fields"
// @Success 200 {object} dtos.EthernetSwitchVLANDto
// @Failure		400		{object}	dtos.ValidationErrorDto
// @Failure		404		"Not Found"
// @Failure		500		"Internal Server Error"
// @router /ethernet-switch/{id}/vlan/{vlanUUID} [put]
func (e *EthernetSwitchVLANGinController) Update(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.EthernetSwitchVLANUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	vlanUUID, err := parseUUIDParam(ctx, "vlanUUID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.UpdateVLAN(ctx, switchID, vlanUUID, reqDto)
	handleWithData(ctx, err, dto)
}

//Delete soft deleting switch VLAN in database
//	Params
//	ctx - gin context
// @Summary Delete ethernet switch VLAN by id
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce	json
// @param	id			path	string		true	"Ethernet switch ID"
// @param	vlanUUID	path	string		true	"Ethernet switch VLAN UUID"
// @Success 204 	"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /ethernet-switch/{id}/vlan/{vlanUUID}  [delete]
func (e *EthernetSwitchVLANGinController) Delete(ctx *gin.Context) {
	switchID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	vlanUUID, err := parseUUIDParam(ctx, "vlanUUID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = e.service.DeleteVLAN(ctx, switchID, vlanUUID)
	handle(ctx, err)
}
