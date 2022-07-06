package controllers

import (
	"net/http"
	"rol/app/interfaces"
	"rol/dtos"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
//Return
//	Returns http status code and response dto
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetList(ctx *gin.Context) {
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
	paginatedList, err := g.service.GetList(ctx, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
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

//GetByID Get entity by id
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}

	dto, err := g.service.GetByID(ctx, id)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
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

//Create new entity
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Create(ctx *gin.Context) {
	reqDto := new(CreateDtoType)
	err := ctx.ShouldBindJSON(&reqDto)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = RestoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	id, err := g.service.Create(ctx, *reqDto)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
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

//Update
//	Update entity in database by id
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Update(ctx *gin.Context) {
	reqDto := new(UpdateDtoType)
	err := ctx.ShouldBindJSON(reqDto)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = RestoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	strID := ctx.Param("id")
	id := uuid.MustParse(strID)

	err = g.service.Update(ctx, *reqDto, id)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Delete
//	Soft deleting entity in database
//Params
//	ctx - gin context
//Return
//	Returns http status code and response dto
func (g *GinGenericController[DtoType, CreateDtoType, UpdateDtoType, EntityType]) Delete(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	err = g.service.Delete(ctx, id)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			g.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}
