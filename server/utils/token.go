package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Invalid Authorization header")
	}

	return parts[1], nil
}
