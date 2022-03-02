package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rol/app/interfaces/generic"
	"rol/dtos"
	"rol/webapi/utils"
	"strconv"
)

type EthernetSwitchController struct {
	service generic.IGenericEntityService
}

func NewEthernetSwitchController(service *generic.IGenericEntityService) *EthernetSwitchController {
	return &EthernetSwitchController{
		service: *service,
	}
}

func (esc *EthernetSwitchController) GetAll(ctx *gin.Context) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	esc.service.GetAll(dtosArr)
	utils.APIResponse(ctx, "Switches was successfully received", http.StatusOK, http.MethodGet, dtosArr)
}

func (esc *EthernetSwitchController) GetById(ctx *gin.Context) {
	dto := dtos.EthernetSwitchDto{}
	strId := ctx.Param("id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {

	}
	id := uint(id64)

	esc.service.GetById(&dto, id)
	if dto.Id > 0 {
		utils.APIResponse(ctx, "The switch was successfully received", http.StatusOK, http.MethodGet, dto)
	} else {
		utils.APIResponse(ctx, "The switch was not found", http.StatusNotFound, http.MethodGet, nil)
	}

}

func (esc *EthernetSwitchController) Create(ctx *gin.Context) {
	dto := dtos.EthernetSwitchCreateDto{}
	err := ctx.ShouldBindJSON(&dto)

	if err != nil {

	}

	esc.service.Create(&dto)
}

func (esc *EthernetSwitchController) Update(ctx *gin.Context) {
	dto := dtos.EthernetSwitchUpdateDto{}
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {

	}

	strId := ctx.Param("id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {

	}
	id := uint(id64)

	esc.service.Update(&dto, id)
}

func (esc *EthernetSwitchController) Delete(ctx *gin.Context) {

}
