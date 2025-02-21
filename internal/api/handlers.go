package api

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func statusHandler(c *gin.Context, cfg *config.Config) {
	c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
}

type CopyRequest struct {
	ID              uuid.UUID `json:"id"`
	SourcePath      string    `json:"source_path"`
	DestinationPath string    `json:"destination_path"`
}

func (r *CopyRequest) EnsureID() {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
}

func copyHandler(c *gin.Context, cfg *config.Config) {
	// Parse the request body
	var req CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.EnsureID()
	// Publish the request to RabbitMQ
	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize request"})
		return
	}

	if err := rabbitmq.Publish(cfg.RabbitMQURL, cfg.FileCopyQueue, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to accept request with id: %s", req.ID.String())})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("File copy request accepted with id: %s", req.ID.String())})
}

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(AuthMiddleware(cfg))
	router.GET("/status", func(c *gin.Context) {
		statusHandler(c, cfg)
	})
	router.POST("/copy", func(c *gin.Context) {
		copyHandler(c, cfg)
	})
	return router
}
