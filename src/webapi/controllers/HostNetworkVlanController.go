package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//HostNetworkVlanController host network vlan API controller
type HostNetworkVlanController struct {
	service *services.HostNetworkVlanService
	logger  *logrus.Logger
}

//NewHostNetworkVlanController host network vlan controller constructor. Parameters pass through DI
//
//Params:
//	vlanService - vlan service
//	log - logrus logger
//Return:
//	*HostNetworkVlanController - instance of host network vlan controller
func NewHostNetworkVlanController(vlanService *services.HostNetworkVlanService, log *logrus.Logger) *HostNetworkVlanController {
	return &HostNetworkVlanController{
		service: vlanService,
		logger:  log,
	}
}

//RegisterHostNetworkVlanController registers controller for getting host network vlans via api
func RegisterHostNetworkVlanController(controller *HostNetworkVlanController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/host/network/vlan/", controller.GetList)
	groupRoute.GET("/host/network/vlan/:name", controller.GetByName)
	groupRoute.POST("/host/network/vlan/", controller.Create)
	groupRoute.PUT("/host/network/vlan/:name", controller.Update)
	groupRoute.DELETE("/host/network/vlan/:name", controller.Delete)
}

//GetList get list of host network VLAN's
//
//Params:
//	ctx - gin context
//
// @Summary Get list of host network VLAN's
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Success	200		{object}	[]dtos.HostNetworkVlanDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/vlan/	[get]
func (h *HostNetworkVlanController) GetList(ctx *gin.Context) {
	vlanList, err := h.service.GetList()
	handleWithData(ctx, err, vlanList)
}

//GetByName get vlan port by name
//
//Params:
//	ctx - gin context
//
// @Summary	Gets vlan port by name
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	name	path		string	true	"Vlan name"
// @Success	200		{object}	dtos.HostNetworkVlanDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/vlan/{name}	[get]
func (h *HostNetworkVlanController) GetByName(ctx *gin.Context) {
	name := ctx.Param("name")
	vlan, err := h.service.GetByName(name)
	handleWithData(ctx, err, vlan)
}

//Create new host vlan
//
//Params:
//	ctx - gin context
//
// @Summary	Create new host vlan
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.HostNetworkVlanCreateDto	true	"Host vlan fields"
// @Success	200		{object}	dtos.HostNetworkVlanDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/vlan/	[post]
func (h *HostNetworkVlanController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDto[dtos.HostNetworkVlanCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	vlanDto, err := h.service.Create(reqDto)
	handleWithData(ctx, err, vlanDto)
}

//Update host network vlan
//
//Params:
//	ctx - gin context
//
// @Summary update host network vlan
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.HostNetworkVlanUpdateDto	true	"Host vlan fields"
// @param	name	path		string							false	"Vlan name"
// @Success	200		{object}	dtos.HostNetworkVlanDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/vlan/{name}	[put]
func (h *HostNetworkVlanController) Update(ctx *gin.Context) {
	reqDto, err := getRequestDto[dtos.HostNetworkVlanUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	name := ctx.Param("name")
	vlanDto, err := h.service.Update(name, reqDto)
	handleWithData(ctx, err, vlanDto)
}

//Delete host network vlan
//
//Params:S
//	ctx - gin context
//
// @Summary	Delete host network vlan by name
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	name	path		string	true	"Vlan name"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/vlan/{name}	[delete]
func (h *HostNetworkVlanController) Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	err := h.service.Delete(name)
	handle(ctx, err)
}
