package controllers

import (
	"github.com/gin-gonic/gin"
	"rol/app/interfaces/generic"
	"rol/domain"
	"rol/dtos"

	"github.com/sirupsen/logrus"
)

//NewEthernetSwitchGinController ethernet switch GIN controller constructor
type EthernetSwitchGinController struct {
	GinGenericController[dtos.EthernetSwitchDto,
		dtos.EthernetSwitchCreateDto,
		dtos.EthernetSwitchUpdateDto,
		domain.EthernetSwitch]
}

// @Summary Gets paginated list of switches
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce json
// @param 	orderBy			path	string	false	"Order by field"
// @param 	orderDirection	path	string	false	"'asc' or 'desc' for ascending or descending order"
// @param 	search 			path	string	false	"searchable value in entity"
// @param 	page 			path	int		false	"page number"
// @param 	pageSize 		path	int		false	"number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.EthernetSwitchDto}} "desc"
// @router /switch/ [get]
func (c *EthernetSwitchGinController) GetList(ctx *gin.Context) {
	c.GinGenericController.GetList(ctx)
}

// @Summary Gets ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @param 	id	path	string		true	"ethernet switch id"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.EthernetSwitchDto}
// @router /switch/{id} [get]
func (c *EthernetSwitchGinController) GetById(ctx *gin.Context) {
	c.GinGenericController.GetById(ctx)
}

// @Summary Creates new ethernet switch
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @Param request body dtos.EthernetSwitchCreateDto true "ethernet switch fields"
// @Success 200 {object} dtos.ResponseDataDto
// @router /switch/ [post]
func (c *EthernetSwitchGinController) Create(ctx *gin.Context) {
	c.GinGenericController.Create(ctx)
}

// @Summary Updates ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @Param request body dtos.EthernetSwitchUpdateDto true "ethernet switch fields"
// @Success 200 {object} dtos.ResponseDto
// @router /switch/{id} [put]
func (c *EthernetSwitchGinController) Update(ctx *gin.Context) {
	c.GinGenericController.Update(ctx)
}

// @Summary Delete ethernet switch by id
// @version 1.0
// @Tags ethernet switch
// @Accept  json
// @Produce  json
// @param 	id	path	string		true	"ethernet switch id"
// @Success 200 {object} dtos.ResponseDto
// @router /switch/{id} [delete]
func (c *EthernetSwitchGinController) Delete(ctx *gin.Context) {
	c.GinGenericController.Delete(ctx)
}

//NewEthernetSwitchGinController ethernet switch controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for ethernet switches
func NewEthernetSwitchGinController(service generic.IGenericService[dtos.EthernetSwitchDto,
	dtos.EthernetSwitchCreateDto,
	dtos.EthernetSwitchUpdateDto,
	domain.EthernetSwitch], log *logrus.Logger) *EthernetSwitchGinController {
	genContr := NewGinGenericController(service, log)
	switchContr := &EthernetSwitchGinController{
		GinGenericController: *genContr,
	}
	return switchContr
}
