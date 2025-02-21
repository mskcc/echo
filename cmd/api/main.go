package main

import (
	"echo/internal/api"
	"echo/internal/config"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the API server
	if err := api.StartServer(cfg); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
