package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//HostNetworkBridgeController host network bridge API controller
type HostNetworkBridgeController struct {
	service *services.HostNetworkService
	logger  *logrus.Logger
}

//NewHostNetworkBridgeController host network bridge controller constructor. Parameters pass through DI
//
//Params:
//	bridgeService - bridge service
//	log - logrus logger
//Return:
//	*HostNetworkBridgeController - instance of host network bridge controller
func NewHostNetworkBridgeController(bridgeService *services.HostNetworkService, log *logrus.Logger) *HostNetworkBridgeController {
	return &HostNetworkBridgeController{
		service: bridgeService,
		logger:  log,
	}
}

//RegisterHostNetworkBridgeController registers controller for getting host network bridges via api
func RegisterHostNetworkBridgeController(controller *HostNetworkBridgeController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/host/network/bridge/", controller.GetList)
	groupRoute.GET("/host/network/bridge/:name", controller.GetByName)
	groupRoute.POST("/host/network/bridge/", controller.Create)
	groupRoute.PUT("/host/network/bridge/:name", controller.Update)
	groupRoute.DELETE("/host/network/bridge/:name", controller.Delete)
}

//GetList get list of host network bridges
//
//Params:
//	ctx - gin context
//
// @Summary Get list of host network bridges
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Success	200		{object}	[]dtos.HostNetworkBridgeDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/bridge/	[get]
func (h *HostNetworkBridgeController) GetList(ctx *gin.Context) {
	bridgeList, err := h.service.GetBridgeList()
	handleWithData(ctx, err, bridgeList)
}

//GetByName get bridge port by name
//
//Params:
//	ctx - gin context
//
// @Summary	Gets bridge port by name
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	name	path		string	true	"Bridge name"
// @Success	200		{object}	dtos.HostNetworkBridgeDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/bridge/{name}	[get]
func (h *HostNetworkBridgeController) GetByName(ctx *gin.Context) {
	name := ctx.Param("name")
	bridge, err := h.service.GetBridgeByName(name)
	handleWithData(ctx, err, bridge)
}

//Create new host bridge
//
//Params:
//	ctx - gin context
//
// @Summary	Create new host bridge
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.HostNetworkBridgeCreateDto	true	"Host bridge fields"
// @Success	200		{object}	dtos.HostNetworkBridgeDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/bridge/	[post]
func (h *HostNetworkBridgeController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.HostNetworkBridgeCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	bridgeDto, err := h.service.CreateBridge(reqDto)
	handleWithData(ctx, err, bridgeDto)
}

//Update host network bridge
//
//Params:
//	ctx - gin context
//
// @Summary update host network bridge
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @Param	name	path		string							true	"Bridge name"
// @Param	request	body		dtos.HostNetworkBridgeUpdateDto	true	"Host bridge fields"
// @Success	200		{object}	dtos.HostNetworkBridgeDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/bridge/{name}	[put]
func (h *HostNetworkBridgeController) Update(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.HostNetworkBridgeUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	name := ctx.Param("name")
	bridgeDto, err := h.service.UpdateBridge(name, reqDto)
	handleWithData(ctx, err, bridgeDto)
}

//Delete host network bridge
//
//Params:
//	ctx - gin context
//
// @Summary	Delete host network bridge by name
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	name	path		string	true	"Bridge name"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/bridge/{name}	[delete]
func (h *HostNetworkBridgeController) Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	err := h.service.DeleteBridge(name)
	handle(ctx, err)
}
