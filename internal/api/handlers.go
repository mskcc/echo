package api

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"echo/internal/worker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func statusHandler(c *gin.Context, cfg *config.Config) {
	c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
}

type CopyRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type DeleteRequest struct {
	Source string `json:"source"`
}

func copyHandler(c *gin.Context, cfg *config.Config) {
	// Parse the request body
	var req CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create FileTask
	fileTask := worker.FileTask{
		ID:          uuid.New(),
		Type:        worker.TaskCopy,
		Source:      req.Source,
		Destination: req.Destination,
	}

	log.Println(fileTask)

	// Publish the request to RabbitMQ
	body, err := json.Marshal(fileTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize request"})
		return
	}

	if err := rabbitmq.Publish(cfg.RabbitMQURL, cfg.TaskQueue, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to accept request with id: %s", fileTask.ID.String())})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("File copy request accepted with id: %s", fileTask.ID.String())})
}

func deleteHandler(c *gin.Context, cfg *config.Config) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create FileTask
	fileTask := worker.FileTask{
		ID:     uuid.New(),
		Type:   worker.TaskDelete,
		Source: req.Source,
	}

	body, err := json.Marshal(fileTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize request"})
		return
	}
	if err := rabbitmq.Publish(cfg.RabbitMQURL, cfg.TaskQueue, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to accept request with id: %s", fileTask.ID.String())})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("Delete file request accepted with id: %s", fileTask.ID.String())})
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
	router.POST("/delete", func(c *gin.Context) {
		deleteHandler(c, cfg)
	})
	return router
}
