package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"rol/app/interfaces"
	"rol/dtos"
)

//GinGenericController generic controller structure for IEntityModel
type GinGenericController[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel] struct {
	//	service - service with needed dtos instantiated
	service interfaces.IGenericService[DtoType, CreateDtoType, UpdateDtoType, EntityType]
	//	logger - logger
	logger *logrus.Logger
}

//NewGinGenericController Constructor for generic controller, works with IEntityDtoModel and IEntityModel interfaces
//Params
//	Instantiate response dto type, create/update dtos and entity to work with controller
//	service - instance of generic service
//	log - logrus logger
//Return
//	*GinGenericController - new generic controller for type which was instantiated
func NewGinGenericController[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType interfaces.IEntityModel](service interfaces.IGenericService[DtoType, CreateDtoType, UpdateDtoType,
	EntityType], log *logrus.Logger) *GinGenericController[DtoType,
	CreateDtoType,
	UpdateDtoType,
	EntityType] {
	return &GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]{
		service: service,
		logger:  log,
	}
}

//GetList Get list of elements with search and pagination
//Params
//	ctx - gin context
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetList(ctx *gin.Context) {
	getListQuery := parseGetListQueryParams(ctx, "id")
	paginatedList, err := g.service.GetList(ctx,
		getListQuery.search,
		getListQuery.orderBy,
		getListQuery.orderDirection,
		getListQuery.page,
		getListQuery.pageSize)
	if err != nil {
		abortByTypedError(ctx, err)
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

//GetByID Get entity by id
//Params
//	ctx - gin context
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetByID(ctx *gin.Context) {
	id, err := parseUUIDFromUrl(ctx, "id")
	if err != nil {
		abortByTypedError(ctx, err)
	}

	dto, err := g.service.GetByID(ctx, id)
	if err != nil {
		abortByTypedError(ctx, err)
	}
	if dto == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
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

//Create new entity
//Params
//	ctx - gin context
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Create(ctx *gin.Context) {
	reqDto := new(CreateDtoType)
	err := ctx.ShouldBindJSON(&reqDto)
	if err != nil {
		abortByTypedError(ctx, err)
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = restoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	id, err := g.service.Create(ctx, *reqDto)
	if err != nil {
		abortByTypedError(ctx, err)
	}
	responseDto := dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: id,
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Update entity by id
//Params
//	ctx - gin context
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Update(ctx *gin.Context) {
	reqDto := new(UpdateDtoType)
	err := ctx.ShouldBindJSON(reqDto)
	if err != nil {
		abortByTypedError(ctx, err)
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = restoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	id, err := parseUUIDFromUrl(ctx, "id")
	if err != nil {
		abortByTypedError(ctx, err)
	}

	err = g.service.Update(ctx, *reqDto, id)
	if err != nil {
		abortByTypedError(ctx, err)
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Delete deleting entity
//Params
//	ctx - gin context
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Delete(ctx *gin.Context) {
	id, err := parseUUIDFromUrl(ctx, "id")
	if err != nil {
		abortByTypedError(ctx, err)
	}

	err = g.service.Delete(ctx, id)
	if err != nil {
		abortByTypedError(ctx, err)
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}
