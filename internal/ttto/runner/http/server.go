package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/qerdcv/ttto/internal/ttto/service"
)

type server struct {
	*gin.Engine

	service *service.Service
}

func newServer(service *service.Service) *server {
	s := &server{
		Engine:  gin.Default(),
		service: service,
	}

	s.Use(s.authMiddleware())

	s.setupRoutes()
	s.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	})
	s.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed",
		})
	})

	return s
}

func (s *server) setupRoutes() {
	g := s.Group("api/v1")
	{
		g.GET("/ping", s.Ping)
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.setupAuthRoutes(g)
}

func (s *server) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
