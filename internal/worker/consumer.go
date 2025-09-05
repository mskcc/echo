package worker

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskCopy   TaskType = "COPY"
	TaskDelete TaskType = "DELETE"
	TaskExists TaskType = "EXISTS"
)

type FileTask struct {
	ID          uuid.UUID `json:"id"`
	Type        TaskType  `json:"type"` // COPY, DELETE, or EXISTS
	Source      string    `json:"source"`
	Destination string    `json:"destination,omitempty"` // Only needed for COPY
}

func (r *FileTask) EnsureID() {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
}

func Start(cfg *config.Config) error {
	// Connect to RabbitMQ
	msgs, err := rabbitmq.Consume(cfg.RabbitMQURL, cfg.TaskQueue)
	if err != nil {
		return err
	}

	log.Println("Worker Service started. Waiting for messages...")

	var wg sync.WaitGroup
	jobs := make(chan FileTask, cfg.NumberOfWorkers) // Buffered channel for worker pool

	// Launch multiple workers
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		wg.Add(1)
		go workerService(uuid.New(), cfg, jobs, &wg)
	}

	// Read messages and send to workers
	for msg := range msgs {
		var req FileTask
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
