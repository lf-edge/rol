// Package controllers describes controllers for webapi
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

//TFTPServerGinController ethernet switch GIN controller constructor
type TFTPServerGinController struct {
	service *services.TFTPServerService
	logger  *logrus.Logger
}

//NewTFTPServerGinController tftp server controller constructor. Parameters pass through DI
//
//Params
//	service - tftp server service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for ethernet switches
func NewTFTPServerGinController(service *services.TFTPServerService, log *logrus.Logger) *TFTPServerGinController {
	tftpContr := &TFTPServerGinController{
		service: service,
		logger:  log,
	}
	return tftpContr
}

//RegisterTFTPServerGinController registers controller for getting ethernet switch ports via api
func RegisterTFTPServerGinController(controller *TFTPServerGinController, server *webapi.GinHTTPServer) {

	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/tftp/", controller.GetList)
	groupRoute.GET("/tftp/:id", controller.GetByID)
	groupRoute.POST("/tftp", controller.Create)
	groupRoute.PUT("/tftp/:id", controller.Update)
	groupRoute.DELETE("/tftp/:id", controller.Delete)

	groupRoute.GET("/tftp/:id/path/", controller.GetPaths)
	groupRoute.POST("/tftp/:id/path/", controller.CreatePath)
	groupRoute.DELETE("/tftp/:id/path/:pathID", controller.DeletePath)
}

//GetList get list of tftp servers with search and pagination
//
//Params
//	ctx - gin context
// @Summary Get paginated list of tftp servers
// @version 1.0
// @Tags	tftp
// @Accept  json
// @Produce json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.TFTPServerDto]
// @Failure	500		"Internal Server Error"
// @router /tftp/ [get]
func (t *TFTPServerGinController) GetList(ctx *gin.Context) {
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
	paginatedList, err := t.service.GetServerList(ctx, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	handleWithData(ctx, err, paginatedList)
}

//GetByID get tftp server by id
//
//Params
//	ctx - gin context
// @Summary	Get TFTP server by id
// @version 1.0
// @Tags	tftp
// @Accept	json
// @Produce	json
// @param	id		path		string				true	"TFTP server ID"
// @Success	200		{object}	dtos.TFTPServerDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id} [get]
func (t *TFTPServerGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := t.service.GetServerByID(ctx, id)
	handleWithData(ctx, err, dto)
}

//Create new TFTP server
//
//Params
//	ctx - gin context
// @Summary	Create new TFTP server
// @version	1.0
// @Tags	tftp
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.TFTPServerCreateDto	true	"TFTP server fields"
// @Success	200		{object}	dtos.TFTPServerDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /tftp/ [post]
func (t *TFTPServerGinController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.TFTPServerCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := t.service.CreateServer(ctx, reqDto)
	handleWithData(ctx, err, dto)
}

//Update TFTP server by id
//
//Params
//	ctx - gin context
// @Summary	Updates TFTP server by id
// @version	1.0
// @Tags	tftp
// @Accept	json
// @Produce	json
// @param	id		path		string					true	"TFTP server ID"
// @Param	request	body		dtos.TFTPServerUpdateDto	true	"TFTP server fields"
// @Success	200		{object}	dtos.TFTPServerDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id} [put]
func (t *TFTPServerGinController) Update(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.TFTPServerUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := t.service.UpdateServer(ctx, reqDto, id)
	handleWithData(ctx, err, dto)
}

//Delete soft deleting TFTP server in database
//
//Params
//	ctx - gin context
// @Summary	Delete TFTP server by id
// @version	1.0
// @Tags	tftp
// @Accept	json
// @Produce	json
// @param	id		path					string		true	"TFTP server ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id} [delete]
func (t *TFTPServerGinController) Delete(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	err = t.service.DeleteServer(ctx, id)
	handle(ctx, err)
}

//GetPaths Get list of TFTP server paths with pagination
//
//Params
//	ctx - gin context
// @Summary	Gets paginated list of TFTP server paths
// @version	1.0
// @Tags	tftp
// @Accept  json
// @Produce	json
// @param	id				path	string		true	"TFTP server ID"
// @param	orderBy			query	string		false	"Order by field"
// @param	orderDirection	query	string		false	"'asc' or 'desc' for ascending or descending order"
// @param	page			query	int			false	"Page number"
// @param	pageSize		query	int			false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.TFTPPathDto]
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id}/path/ [get]
func (t *TFTPServerGinController) GetPaths(ctx *gin.Context) {
	orderBy := ctx.DefaultQuery("orderBy", "ActualPath")
	orderDirection := ctx.DefaultQuery("orderDirection", "asc")
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
	serverID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := t.service.GetPathsList(ctx, serverID, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	handleWithData(ctx, err, paginatedList)
}

//CreatePath new TFTP server path
//
//Params
//	ctx - gin context
// @Summary Creates new TFTP server path
// @version 1.0
// @Tags tftp
// @Accept	json
// @Produce	json
// @param	id		path		string						true	"TFTP server ID"
// @Param	request	body		dtos.TFTPPathCreateDto		true    "TFTP server path fields"
// @Success	200		{object}	dtos.TFTPPathDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id}/path/ [post]
func (t *TFTPServerGinController) CreatePath(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.TFTPPathCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	serverID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := t.service.CreatePath(ctx, serverID, reqDto)
	handleWithData(ctx, err, dto)
}

//DeletePath deleting TFTP server path
//
//Params
//	ctx - gin context
// @Summary Delete TFTP server path by id
// @version 1.0
// @Tags tftp
// @Accept	json
// @Produce	json
// @param	id		path					string	true	"TFTP server ID"
// @param	pathID	path					string	true	"TFTP server path ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /tftp/{id}/path/{pathID} [delete]
func (t *TFTPServerGinController) DeletePath(ctx *gin.Context) {
	serverID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	pathID, err := parseUUIDParam(ctx, "pathID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = t.service.DeletePath(ctx, serverID, pathID)
	handle(ctx, err)
}
