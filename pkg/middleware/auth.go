package middleware

import (
	"net/http"

	utils "api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header value
		authHeader := c.GetHeader("Authorization")
		err := utils.ValidToken(authHeader)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Next()
	}
}
