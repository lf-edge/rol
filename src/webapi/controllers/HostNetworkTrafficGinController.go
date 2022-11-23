// Package controllers describes controllers for webapi
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//HostNetworkTrafficGinController host network bridge API controller
type HostNetworkTrafficGinController struct {
	service *services.HostNetworkService
	logger  *logrus.Logger
}

//NewHostNetworkTrafficGinController host network traffic controller constructor. Parameters pass through DI
//
//Params:
//	trafficService - bridge service
//	log - logrus logger
//Return:
//	*HostNetworkBridgeController - instance of host network bridge controller
func NewHostNetworkTrafficGinController(trafficService *services.HostNetworkService, log *logrus.Logger) *HostNetworkTrafficGinController {
	return &HostNetworkTrafficGinController{
		service: trafficService,
		logger:  log,
	}
}

//RegisterHostNetworkTrafficGinController registers controller for getting host network bridges via api
func RegisterHostNetworkTrafficGinController(controller *HostNetworkTrafficGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/host/network/traffic/:table", controller.GetTableRules)
	groupRoute.POST("/host/network/traffic/:table/", controller.Create)
	groupRoute.DELETE("/host/network/traffic/:table/", controller.Delete)
}

//GetTableRules get selected netfilter table rules
//
//Params:
//	ctx - gin context
//
// @Summary Get selected netfilter table rules
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	table	path		string	true	"Table name"
// @Success	200		{object}	[]dtos.HostNetworkTrafficRuleDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/traffic/{table}	[get]
func (h *HostNetworkTrafficGinController) GetTableRules(ctx *gin.Context) {
	table := ctx.Param("table")
	rules, err := h.service.GetTableRules(table)
	handleWithData(ctx, err, rules)
}

//Create new traffic rule in specified table
//
//Params:
//	ctx - gin context
//
// @Summary	Create new traffic rule in specified table
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	table	path		string	true	"Table name"
// @param	request	body		dtos.HostNetworkTrafficRuleCreateDto	true	"Host traffic rule fields"
// @Success	200		{object}	dtos.HostNetworkTrafficRuleDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router	/host/network/traffic/{table}	[post]
func (h *HostNetworkTrafficGinController) Create(ctx *gin.Context) {
	table := ctx.Param("table")
	reqDto, err := getRequestDtoAndRestoreBody[dtos.HostNetworkTrafficRuleCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	ruleDto, err := h.service.CreateTrafficRule(table, reqDto)
	handleWithData(ctx, err, ruleDto)
}

//Delete netfilter traffic rule in specified table
//
//Params:
//	ctx - gin context
//
// @Summary	Delete netfilter traffic rule in specified table
// @version	1.0
// @Tags	host
// @Accept	json
// @Produce	json
// @param	table	path		string	true	"Table name"
// @param	request	body		dtos.HostNetworkTrafficRuleDeleteDto	true	"Host traffic rule fields"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/host/network/traffic/{table}	[delete]
func (h *HostNetworkTrafficGinController) Delete(ctx *gin.Context) {
	table := ctx.Param("table")
	reqDto, err := getRequestDtoAndRestoreBody[dtos.HostNetworkTrafficRuleDeleteDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = h.service.DeleteTrafficRule(table, reqDto)
	handle(ctx, err)
}
