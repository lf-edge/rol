package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//DeviceGinController Device entity GIN controller constructor
type DeviceGinController struct {
	service *services.DeviceService
	logger  *logrus.Logger
}

//NewDeviceGinController device entity controller constructor
//Params
//	service - device entity service
//	log - logrus logger
//Return
//	*DeviceGinController - Gin controller for Device entity
func NewDeviceGinController(service *services.DeviceService, log *logrus.Logger) *DeviceGinController {
	deviceContr := &DeviceGinController{
		service: service,
		logger:  log,
	}
	return deviceContr
}

//RegisterDeviceGinController registers controller for the Device entity
func RegisterDeviceGinController(controller *DeviceGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/device/", controller.GetDeviceList)
	groupRoute.GET("/device/:id", controller.GetDeviceByID)
	groupRoute.POST("/device", controller.CreateDevice)
	groupRoute.PUT("/device/:id", controller.UpdateDevice)
	groupRoute.DELETE("/device/:id", controller.DeleteDevice)
	//deviceNetInterface
	groupRoute.GET("/device/:id/network-interface", controller.GetDeviceNetInterfacesList)
	groupRoute.GET("/device/:id/network-interface/:interfaceID", controller.GetDeviceNetInterfaceByID)
	groupRoute.POST("/device/:id/network-interface", controller.CreateDeviceNetInterface)
	groupRoute.PUT("/device/:id/network-interface/:interfaceID", controller.UpdateDeviceNetInterface)
	groupRoute.DELETE("/device/:id/network-interface/:interfaceID", controller.DeleteDeviceNetInterface)
	//PowerCommand
	groupRoute.POST("/device/:id/power-command", controller.PowerCommand)
}

//GetDeviceList Get list of devices with search and pagination
//Params
//	ctx - gin context
// @Summary Gets paginated list of device entity
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"searchable value in entity"
// @param	page			query	int		false	"page number"
// @param	pageSize		query	int		false	"number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.DeviceDto]
// @Failure	500		"Internal Server Error"
// @router /device/ [get]
func (c *DeviceGinController) GetDeviceList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "CreatedAt", "asc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := c.service.GetDevicesList(ctx, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetDeviceByID get device by id
//	Params
//	ctx - gin context
// @Summary	Get device by id
// @version 1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Device ID"
// @Success	200		{object}	dtos.DeviceDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /device/{id} [get]
func (c *DeviceGinController) GetDeviceByID(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := c.service.GetDeviceByID(ctx, id)
	handleWithData(ctx, err, dto)
}

//CreateDevice - create new device
//	Params
//	ctx - gin context
// @Summary	Create device
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.DeviceCreateDto	true	"Device fields"
// @Success	200		{object}	dtos.DeviceDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /device/ [post]
func (c *DeviceGinController) CreateDevice(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DeviceCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := c.service.CreateDevice(ctx, reqDto)
	handleWithData(ctx, err, dto)
}

//UpdateDevice - update device
//	Params
//	ctx - gin context
// @Summary	Updates device
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Device ID"
// @Param	request	body		dtos.DeviceUpdateDto true "Device fields"
// @Success	200		{object}	dtos.DeviceDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /device/{id} [put]
func (c *DeviceGinController) UpdateDevice(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DeviceUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := c.service.UpdateDevice(ctx, id, reqDto)
	handleWithData(ctx, err, dto)
}

//DeleteDevice - delete device by id
//	Params
//	ctx - gin context
// @Summary	Delete device by id
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id		path	string		true	"Device ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /device/{id} [delete]
func (c *DeviceGinController) DeleteDevice(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	err = c.service.DeleteDevice(ctx, id)
	handle(ctx, err)
}

//GetDeviceNetInterfacesList - get list of device network interfaces with search and pagination
//	Params
//	ctx - gin context
// @Summary Get list of device network interface
// @version 1.0
// @Tags	device
// @Accept  json
// @Produce json
// @param	id				path	string	true	"Device ID"
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.DeviceNetworkInterfaceDto]
// @Failure	500		"Internal Server Error"
// @router /device/{id}/network-interface [get]
func (c *DeviceGinController) GetDeviceNetInterfacesList(ctx *gin.Context) {
	req := newPaginatedRequestStructForParsing(1, 10, "CreatedAt", "asc", "")
	err := parseGinRequest(ctx, &req)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	paginatedList, err := c.service.GetNetInterfacesList(ctx, id, req.Search, req.OrderBy, req.OrderDirection,
		req.Page, req.PageSize)
	handleWithData(ctx, err, paginatedList)
}

//GetDeviceNetInterfaceByID get device network interface by id
//	Params
//	ctx - gin context
// @Summary	Get device network interface by id
// @version 1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id			path		string		true	"Device ID"
// @param	interfaceID	path		string		true	"Device network interface ID"
// @Success	200			{object}	dtos.DeviceNetworkInterfaceDto
// @Failure	404			"Not Found"
// @Failure	500			"Internal Server Error"
// @router /device/{id}/network-interface/{interfaceID} [get]
func (c *DeviceGinController) GetDeviceNetInterfaceByID(ctx *gin.Context) {
	deviceID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	interfaceID, err := parseUUIDParam(ctx, "interfaceID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	dto, err := c.service.GetNetInterfaceByID(ctx, deviceID, interfaceID)
	handleWithData(ctx, err, dto)
}

//CreateDeviceNetInterface - crate new device network interface
//	Params
//	ctx - gin context
// @Summary	Create device network interface
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"DHCP v4 server ID"
// @Param	request	body		dtos.DeviceNetworkInterfaceCreateDto	true	"Device network interface fields"
// @Success	200		{object}	dtos.DeviceNetworkInterfaceDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router /device/{id}/network-interface [post]
func (c *DeviceGinController) CreateDeviceNetInterface(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DeviceNetworkInterfaceCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := c.service.CreateNetInterface(ctx, id, reqDto)
	handleWithData(ctx, err, dto)
}

//UpdateDeviceNetInterface - update device network interface
//	Params
//	ctx - gin context
// @Summary	Updates device network interface
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id			path		string		true	"Device ID"
// @param	interfaceID	path		string		true	"Device network interface ID"
// @Param	request		body		dtos.DeviceNetworkInterfaceUpdateDto true "Device network interface fields"
// @Success	200			{object}	dtos.DeviceNetworkInterfaceDto
// @Failure	400			{object}	dtos.ValidationErrorDto
// @Failure	404			"Not Found"
// @Failure	500			"Internal Server Error"
// @router /device/{id}/network-interface/{interfaceID} [put]
func (c *DeviceGinController) UpdateDeviceNetInterface(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DeviceNetworkInterfaceUpdateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	serverID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	interfaceID, err := parseUUIDParam(ctx, "interfaceID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	dto, err := c.service.UpdateNetInterface(ctx, serverID, interfaceID, reqDto)
	handleWithData(ctx, err, dto)
}

//DeleteDeviceNetInterface - delete device network interface
//	Params
//	ctx - gin context
// @Summary	Delete device network interface by id
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id			path		string		true	"Device ID"
// @param	interfaceID	path		string		true	"Device network interface ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /device/{id}/network-interface/{interfaceID} [delete]
func (c *DeviceGinController) DeleteDeviceNetInterface(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	interfaceID, err := parseUUIDParam(ctx, "interfaceID")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	err = c.service.DeleteNetInterface(ctx, id, interfaceID)
	handle(ctx, err)
}

//PowerCommand - send power device command
//	Params
//	ctx - gin context
// @Summary	Send power device command
// @version	1.0
// @Tags	device
// @Accept	json
// @Produce	json
// @param	id		path		string		true	"Device ID"
// @Param	request	body		dtos.DevicePowerCommandDto true "Device power command fields"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router /device/{id}/power-command/ [post]
func (c *DeviceGinController) PowerCommand(ctx *gin.Context) {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	reqDto, err := getRequestDtoAndRestoreBody[dtos.DevicePowerCommandDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = c.service.DevicePowerCommand(ctx, id, reqDto)
	handle(ctx, err)
}
