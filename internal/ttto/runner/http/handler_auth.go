package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

func (s *server) setupAuthRoutes(g *gin.RouterGroup) {
	authG := g.Group("/auth")
	{
		authG.POST("/register", s.Register)
		authG.POST("/login", s.Login)
		authG.GET("/logout", s.Register)
	}

}

// Register godoc
//
//	@Summary		register new user
//	@Description	creates new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		http.UserRequest	true	"Data of new user"
//	@Success		201		{object}	http.Response
//	@Failure		400		{object}	http.Response
//	@Failure		409		{object}	http.Response
//	@Failure		422		{object}	http.Response
//	@Failure		500		{object}	http.Response
//	@Router			/api/v1/auth/register [post]
func (s *server) Register(c *gin.Context) {
	var u UserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if err := s.service.CreateUser(c.Request.Context(), domain.User{
		Username: u.Username,
		Password: u.Password,
	}); err != nil {
		var valErr *service.ErrValidation
		switch {
		case errors.As(err, &valErr):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": valErr,
			})
		case errors.Is(err, service.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"message": err,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": http.StatusText(http.StatusCreated),
	})
}

// Login godoc
//
//	@Summary		login to the app
//	@Description	logins to the application
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		http.UserRequest	true	"Data of new user"
//	@Success		201		{object}	http.Response
//	@Failure		400		{object}	http.Response
//	@Failure		409		{object}	http.Response
//	@Failure		422		{object}	http.Response
//	@Failure		500		{object}	http.Response
//	@Router			/api/v1/auth/login [post]
func (s *server) Login(c *gin.Context) {
	var u UserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	token, err := s.service.AuthorizeUser(c.Request.Context(), domain.User{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		var valErr *service.ErrValidation
		switch {
		case errors.As(err, &valErr):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": valErr,
			})
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": service.ErrInvalidCredentials,
			})
		default:
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	c.SetCookie("token", token, int(time.Hour*24*7), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}

func (s *server) Logout(c *gin.Context) {
}
