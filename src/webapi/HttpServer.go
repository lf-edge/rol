package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rol/app/interfaces/generic"
	"rol/webapi/controllers"
	"rol/webapi/middleware"
)

type HttpServer struct {
	engine  *gin.Engine
	service *generic.IGenericEntityService
}

func NewHttpServer(logger *logrus.Logger, service *generic.IGenericEntityService) HttpServer {
	ginEngine := gin.New()
	ginEngine.Use(middleware.Logger(logger), gin.Recovery())
	serv := HttpServer{
		engine:  ginEngine,
		service: service,
	}
	return serv
}

func (server *HttpServer) InitializeRoutes() {
	server.InitializeControllers()
}

func (server *HttpServer) InitializeControllers() error {
	switchContr := controllers.NewEthernetSwitchController(server.service)

	groupRoute := server.engine.Group("/api/v1")
	groupRoute.GET("/switch/:id", switchContr.GetById)
	groupRoute.GET("/switch", switchContr.GetAll)
	groupRoute.POST("/switch", switchContr.Create)
	groupRoute.PUT("/switch/:id", switchContr.Update)

	return nil
}

func (server *HttpServer) Start() {
	server.InitializeRoutes()
	server.engine.Run("localhost:8080")
}
