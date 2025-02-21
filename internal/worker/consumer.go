package worker

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"sync"
)

type CopyFileRequest struct {
	ID              uuid.UUID `json:"id"`
	SourcePath      string    `json:"source_path"`
	DestinationPath string    `json:"destination_path"`
}

func (r *CopyFileRequest) EnsureID() {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
}

func Start(cfg *config.Config) error {
	// Connect to RabbitMQ
	msgs, err := rabbitmq.Consume(cfg.RabbitMQURL, cfg.FileCopyQueue)
	if err != nil {
		return err
	}

	log.Println("Worker Service started. Waiting for messages...")

	var wg sync.WaitGroup
	jobs := make(chan CopyFileRequest, cfg.NumberOfWorkers) // Buffered channel for worker pool

	// Launch multiple workers
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		wg.Add(1)
		go copyWorkerService(uuid.New(), cfg, jobs, &wg)
	}

	// Read messages and send to workers
	for msg := range msgs {
		var req CopyFileRequest
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}
		req.EnsureID()
		jobs <- req
	}

	close(jobs)
	wg.Wait()

	return nil
}
