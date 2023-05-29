package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

func (s *server) setupAuthRoutes(g *gin.RouterGroup) {
	authG := g.Group("")
	{
		authG.POST("/registration", s.registration)
		authG.POST("/login", s.login)
		authG.GET("/logout", s.logout)
	}

}

// register godoc
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
//	@Router			/api/registration [post]
func (s *server) registration(c *gin.Context) {
	var u UserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if err := s.service.CreateUser(c.Request.Context(), &domain.User{
		Username: u.Username,
		Password: u.Password,
	}); err != nil {
		var valErr *service.ErrValidation
		switch {
		case errors.As(err, &valErr):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": valErr,
			})
		case errors.Is(err, service.ErrUserAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": http.StatusText(http.StatusConflict),
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": http.StatusText(http.StatusCreated),
	})
}

// login godoc
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
//	@Router			/api/login [post]
func (s *server) login(c *gin.Context) {
	var u UserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	token, err := s.service.AuthorizeUser(c.Request.Context(), &domain.User{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		var valErr *service.ErrValidation
		switch {
		case errors.As(err, &valErr):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": valErr,
			})
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrInvalidCredentials):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": service.ErrInvalidCredentials.Error(),
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	c.SetCookie("token", token, int((time.Hour * 24 * 7).Milliseconds()), "/", "*", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}

// login godoc
//
//	@Summary		logout from the app
//	@Description	logouts from the application
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	http.Response
//	@Failure		500	{object}	http.Response
//	@Router			/api/logout [get]
func (s *server) logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "*", false, true)
}
