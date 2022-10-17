package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//DHCP4ServerGinController DHCP v4 server GIN controller constructor
type DHCP4ServerGinController struct {
	service *services.DHCP4ServerService
	logger  *logrus.Logger
}

//RegisterDHCP4ServerGinController registers controller for the DHCP v4 servers
func RegisterDHCP4ServerGinController(controller *DHCP4ServerGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/dhcp/", controller.GetServersList)
	groupRoute.GET("/dhcp/:id", controller.GetServerByID)
	groupRoute.POST("/dhcp", controller.CreateServer)
	groupRoute.PUT("/dhcp/:id", controller.UpdateServer)
	groupRoute.DELETE("/dhcp/:id", controller.DeleteServer)
}

//NewDHCP4ServerGinController dhcp v4 server controller constructor. Parameters pass through DI
//Params
//	service - dhcp v4 server service
//	log - logrus logger
//Return
//	*DHCP4ServerGinController - Gin controller for dhcp v4 servers
func NewDHCP4ServerGinController(service *services.DHCP4ServerService, log *logrus.Logger) *DHCP4ServerGinController {
	switchContr := &DHCP4ServerGinController{
		service: service,
		logger:  log,
	}
	return switchContr
}

//GetServersList get list of dhcp v4 servers with search and pagination
//	Params
//	ctx - gin context
// @Summary Get paginated list of dhcp v4 servers
// @version 1.0
// @Tags	dhcp
// @Accept  json
// @Produce json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.DHCP4ServerDto]
// @Failure	500		"Internal Server Error"
// @router /dhcp/ [get]
func (e *DHCP4ServerGinController) GetServersList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "CreatedAt", "asc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := e.service.GetServerList(ctx, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetServerByID get dhcp v4 server by id
//	Params
//	ctx - gin context
// @Summary	Get dhcp v4 server by id
// @version 1.0
// @Tags	dhcp
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"DHCP v4 server ID"
// @Success	200		{object}	dtos.DHCP4ServerDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /dhcp/{id} [get]
func (e *DHCP4ServerGinController) GetServerByID(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := e.service.GetServerByID(ctx, id)
	handleWithData(ctx, err, dto)
}

//CreateServer new DHCP v4 server
//	Params
//	ctx - gin context
// @Summary	Create DHCP v4 server
// @version	1.0
// @Tags	dhcp
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.DHCP4ServerCreateDto	true	"DHCP v4 server fields"
// @Success	200		{object}	dtos.DHCP4ServerDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /dhcp/ [post]
func (e *DHCP4ServerGinController) CreateServer(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DHCP4ServerCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := e.service.CreateServer(ctx, reqDto)
	handleWithData(ctx, err, dto)
}

//UpdateServer DHCP v4 server by id
//	Params
//	ctx - gin context
// @Summary	Updates DHCP v4 server by id
// @version	1.0
// @Tags	dhcp
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"DHCP v4 server ID"
// @Param	request	body		dtos.DHCP4ServerUpdateDto true "DHCP v4 server fields"
// @Success	200		{object}	dtos.DHCP4ServerDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /dhcp/{id} [put]
func (e *DHCP4ServerGinController) UpdateServer(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DHCP4ServerUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := e.service.UpdateServer(ctx, id, reqDto)
	handleWithData(ctx, err, dto)
}

//DeleteServer deleting dhcp v4 server
//	Params
//	ctx - gin context
// @Summary	Delete dhcp v4 server by id
// @version	1.0
// @Tags	dhcp
// @Accept	json
// @Produce	json
// @param	id		path	string		true	"DHCP v4 server ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /dhcp/{id} [delete]
func (e *DHCP4ServerGinController) DeleteServer(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	err = e.service.DeleteServer(ctx, id)
	handle(ctx, err)
}
