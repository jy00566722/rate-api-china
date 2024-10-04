package middleware

import (
	"noFree/config"
	"noFree/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := parts[1]
		claims, err := utils.ValidateJWT(cfg, token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func DeviceAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fingerprint := c.GetHeader("X-Device-Fingerprint")
		if fingerprint == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Device fingerprint is required"})
			return
		}

		c.Set("fingerprint", fingerprint)
		c.Next()
	}
}
