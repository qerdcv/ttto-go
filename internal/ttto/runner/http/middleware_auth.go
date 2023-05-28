package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qerdcv/ttto/internal/xctx"
)

func (s *server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err == nil {
			user, decErr := s.service.DecodeToken(cookie)
			if decErr == nil {
				c.Request = c.Request.WithContext(xctx.ContextWithUser(c.Request.Context(), user))
			}
		}

		c.Next()
	}
}

func (s *server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := xctx.UserFromContext(c.Request.Context())
		if u == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		c.Next()
	}
}
