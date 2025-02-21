package main

import (
	"log"

	"echo/internal/config"
	"echo/internal/worker"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the consumer
	if err := worker.Start(cfg); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
