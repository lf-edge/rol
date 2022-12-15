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

//ProjectGinController user project API controller
type ProjectGinController struct {
	service *services.ProjectService
	logger  *logrus.Logger
}

//NewProjectGinController user project controller constructor. Parameters pass through DI
//
//Params:
//	projectService - user project service
//	log - logrus logger
//Return:
//	*ProjectGinController - instance of user project controller
func NewProjectGinController(projectService *services.ProjectService, log *logrus.Logger) *ProjectGinController {
	return &ProjectGinController{
		service: projectService,
		logger:  log,
	}
}

//RegisterProjectGinController registers controller for getting user projects via api
func RegisterProjectGinController(controller *ProjectGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/project/", controller.GetList)
	groupRoute.GET("/project/:id", controller.GetByID)
	groupRoute.POST("/project/", controller.Create)
	groupRoute.DELETE("/project/:id", controller.Delete)
}

//GetList get list of user projects
//
//Params:
//	ctx - gin context
//
// @Summary Get list of user projects
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @param	orderBy			query	string	false	"Order by field"
// @param	orderDirection	query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	search			query	string	false	"Searchable value in entity"
// @param	page			query	int		false	"Page number"
// @param	pageSize		query	int		false	"Number of entities per page"
// @Success	200		{object}	dtos.PaginatedItemsDto[dtos.ProjectDto]
// @Failure	500		"Internal Server Error"
// @router	/project/	[get]
func (p *ProjectGinController) GetList(ctx *gin.Context) {
	orderBy := ctx.DefaultQuery("orderBy", "id")
	orderDirection := ctx.DefaultQuery("orderDirection", "asc")
	search := ctx.DefaultQuery("search", "")
	page := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	projList, err := p.service.GetList(ctx, search, orderBy, orderDirection, pageInt, pageSizeInt)
	handleWithData(ctx, err, projList)
}

//GetByID get user project by id
//
//Params:
//	ctx - gin context
//
// @Summary	Gets user project by id
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @param	id	path		string	true	"User project ID"
// @Success	200		{object}	dtos.ProjectDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/project/{id}	[get]
func (p *ProjectGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	vlan, err := p.service.GetByID(ctx, id)
	handleWithData(ctx, err, vlan)
}

//Create new user project
//
//Params:
//	ctx - gin context
//
// @Summary	Create new user project
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.ProjectCreateDto	true	"User project fields"
// @Success	200		{object}	dtos.ProjectDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router	/project/	[post]
func (p *ProjectGinController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.ProjectCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	projectDto, err := p.service.Create(ctx, reqDto)
	handleWithData(ctx, err, projectDto)
}

//Delete user project
//
//Params:
//	ctx - gin context
//
// @Summary	Delete user project by id
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @param	id	path		string	true	"User project ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/project/{id}	[delete]
func (p *ProjectGinController) Delete(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = p.service.Delete(ctx, id)
	handle(ctx, err)
}
