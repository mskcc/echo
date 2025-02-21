package worker

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type CopyFileResponse struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
}

func copyWorkerService(id uuid.UUID, cfg *config.Config, jobs <-chan CopyFileRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %s started", id.String())
	for req := range jobs {
		log.Printf("Worker %s is processing file copy: %s -> %s", id, req.SourcePath, req.DestinationPath)

		var status string
		var message string
		if err := copyFile(req.SourcePath, req.DestinationPath); err != nil {
			log.Printf("Failed to copy file: %v", err)
			status = "fail"
			message = fmt.Sprintf("Failed to copy file: %s -> %s", req.SourcePath, req.DestinationPath)
		} else {
			log.Printf("File copied successfully: %s -> %s", req.SourcePath, req.DestinationPath)
			status = "success"
			message = fmt.Sprintf("File copied successfully: %s -> %s", req.SourcePath, req.DestinationPath)
		}

		// Publish response back to RabbitMQ
		msg := CopyFileResponse{
			ID:      req.ID,
			Status:  status,
			Message: message,
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to serialize request body: %v", err)
			continue
		}
		rabbitmq.Publish(cfg.RabbitMQURL, cfg.ConfirmationQueue, body)
		log.Printf("Worker %s completed processing job %s", id.String(), req.ID.String())
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination directory if it doesn't exist
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
