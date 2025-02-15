package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PrivilegeMiddleware(requiredLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		privilege, exists := c.Get("privilege")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Privilege level not found"})
			return
		}

		if privilege.(int) < requiredLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
			return
		}

		c.Next()
	}
}
