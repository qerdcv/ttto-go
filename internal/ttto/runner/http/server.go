package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/eventst"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

type server struct {
	*gin.Engine

	es      *eventst.EventStream[*domain.Game]
	service *service.Service
}

func newServer(service *service.Service, es *eventst.EventStream[*domain.Game]) *server {
	s := &server{
		Engine:  gin.Default(),
		service: service,
		es:      es,
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
	g := s.Group("api")
	{
		g.GET("/ping", s.ping)
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.setupAuthRoutes(g)
	s.setupGameHandlers(g)
}

// ping godoc
//
//	@Summary		pong
//	@Description	pong
//	@Tags			healthcheck
//	@Produce		json
//	@Success		200	{object}	http.Response
//	@Failure		500	{object}	http.Response
//	@Router			/api/v1/ping [get]
func (s *server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
