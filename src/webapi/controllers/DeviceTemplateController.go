package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
	"strconv"
)

//DeviceTemplateController device template controller structure
type DeviceTemplateController struct {
	service *services.DeviceTemplateService
	logger  *logrus.Logger
}

//NewDeviceTemplateController constructor for device template controller
//Params
//	service - instance of device template service
//	log - logrus logger
//Return
//	*DeviceTemplateController - new device template controller
//	error - if an error occurred, otherwise nil
func NewDeviceTemplateController(service *services.DeviceTemplateService, log *logrus.Logger) (*DeviceTemplateController, error) {
	return &DeviceTemplateController{
		service: service,
		logger:  log,
	}, nil
}

//RegisterDeviceTemplateController registers controller for getting device templates via api on path /template/device/
func RegisterDeviceTemplateController(controller *DeviceTemplateController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/template/device/", controller.GetList)
	groupRoute.GET("/template/device/:name", controller.GetByName)
}

//GetList Get list of elements with search and pagination
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
// @Summary Gets paginated list of device templates
// @version 1.0
// @Tags device template
// @Accept  json
// @Produce  json
// @param 	orderBy			path	string	false	"Order by field"
// @param 	orderDirection	path	string	false	"'asc' or 'desc' for ascending or descending order"
// @param 	search 			path	string	false	"searchable value in entity"
// @param 	page 			path	int		false	"page number"
// @param 	pageSize 		path	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.DeviceTemplateDto}} ""
// @router /template/device/ [get]
func (d *DeviceTemplateController) GetList(ctx *gin.Context) {
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
	paginatedList, err := d.service.GetList(ctx, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			d.logger.Errorf("%s : %s", err, controllerErr)
		}
	}
	responseDto := &dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: paginatedList,
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//GetByName Get device template by name
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
// @Summary Gets device template by its name
// @version 1.0
// @Tags device template
// @Accept  json
// @Produce  json
// @param 	name	path	string		true	"device template name"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.DeviceTemplateDto}
// @router /template/device/{name} [get]
func (d *DeviceTemplateController) GetByName(ctx *gin.Context) {
	name := ctx.Param("name")

	dto, err := d.service.GetByName(ctx, name)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			d.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}
	if dto == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	responseDto := &dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: dto,
	}
	ctx.JSON(http.StatusOK, responseDto)
}
