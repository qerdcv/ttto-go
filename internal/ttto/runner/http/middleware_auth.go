package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err == nil {
			user, err := s.service.DecodeToken(cookie)
			fmt.Println(err)
			if err == nil {
				fmt.Println("token user", user)
				c.Set("user", user)
			}
		}

		fmt.Println(err)
		c.Next()
	}
}
