package webapi

import (
	"fmt"
	"rol/domain"
	"rol/webapi/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//GinHTTPServer Gin HTTP server struct
type GinHTTPServer struct {
	Engine  *gin.Engine
	logger  *logrus.Logger
	address string
}

//NewGinHTTPServer gin server constructor
//Params
//	log - logrus logger
//	config - application configurations
//Return
//	*GinHTTPServer - gin http server instance
func NewGinHTTPServer(log *logrus.Logger, config *domain.AppConfig) *GinHTTPServer {
	ginEngine := gin.New()
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	ginEngine.Use(middleware.Logger(log), middleware.Recovery(log))
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return &GinHTTPServer{
		Engine:  ginEngine,
		logger:  log,
		address: fmt.Sprintf("%s:%s", config.HTTPServer.Host, config.HTTPServer.Port),
	}
}

//Start starts http server
func (g *GinHTTPServer) Start() {
	err := g.Engine.Run(g.address)
	if err != nil {
		g.logger.Errorf("[HTTP server] start server error: %s", err.Error())
		return
	}
}

//StartHTTPServer starts a new http server from fx.Invoke
func StartHTTPServer(server *GinHTTPServer) {
	server.Start()
}
