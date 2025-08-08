package api

import (
	"log"

	"echo/internal/config"
)

// StartServer starts the Gin HTTP server.
func StartServer(cfg *config.Config) error {
	// Set up the Gin router
	router := SetupRouter(cfg)
	// Start the server
	log.Println("API Service started on :" + cfg.Port)
	return router.Run(":" + cfg.Port)
}
