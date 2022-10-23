package controllers

import (
	"github.com/gin-gonic/gin"
	"rol/app/services"
	"rol/webapi"
)

//HostNetworkController host network API controller
type HostNetworkController struct {
	service *services.HostNetworkService
}

//NewHostNetworkController host network controller constructor
func NewHostNetworkController(service *services.HostNetworkService) *HostNetworkController {
	return &HostNetworkController{service: service}
}

//RegisterHostNetworkController registers controller for management of host network settings
func RegisterHostNetworkController(controller *HostNetworkController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/host/network/ping", controller.Ping)
}

//Ping calls the backend to notify that the current setting does not break the connection
//
//Params:
//	ctx - gin context
//
// @Summary Call the backend to notify that the current setting does not break the connection
// @version	1.0
// @Tags	host
// @Success	204
// @Failure	500		"Internal Server Error"
// @router	/host/network/ping	[get]
func (c *HostNetworkController) Ping(ctx *gin.Context) {
	err := c.service.Ping()
	handle(ctx, err)
}
