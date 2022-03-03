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

func (esc *EthernetSwitchController) GetList(ctx *gin.Context) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
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

	esc.service.GetList(dtosArr, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	utils.APIResponse(ctx, "Switches was successfully received", http.StatusOK, http.MethodGet, dtosArr)
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
		ctx.AbortWithError(http.StatusBadRequest, err)
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
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	esc.service.Create(&dto)
}

func (esc *EthernetSwitchController) Update(ctx *gin.Context) {
	dto := dtos.EthernetSwitchUpdateDto{}
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	strId := ctx.Param("id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	id := uint(id64)

	esc.service.Update(&dto, id)
}

func (esc *EthernetSwitchController) Delete(ctx *gin.Context) {

}
