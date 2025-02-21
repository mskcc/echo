package api

import (
	"echo/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware is a middleware that checks for a valid API token.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer "+cfg.APIToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing
			return
		}
		c.Next()
	}
}
