package worker

import (
	"echo/internal/config"
	"echo/internal/rabbitmq"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type TaskResponse struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
}

func workerService(id uuid.UUID, cfg *config.Config, jobs <-chan FileTask, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %s started", id.String())
	for req := range jobs {
		log.Printf("Worker %s is processing task type %s", id, req.Type)
		var status string
		var message string
		switch req.Type {

		case TaskCopy:

			if err := copyFile(req.Source, req.Destination); err != nil {
				log.Printf("Failed to copy file: %v", err)
				status = "fail"
				message = fmt.Sprintf("Failed to copy file: %s -> %s", req.Source, req.Destination)
			} else {
				log.Printf("File copied successfully: %s -> %s", req.Source, req.Destination)
				status = "success"
				message = fmt.Sprintf("File copied successfully: %s -> %s", req.Source, req.Destination)
			}
		case TaskDelete:
			if err := deleteFile(req.Source); err != nil {
				log.Printf("Failed to delete file: %v", err)
				status = "fail"
				message = fmt.Sprintf("Failed to delete file: %s", req.Source)
			} else {
				log.Printf("File deleted successfully: %s", req.Source)
				status = "success"
				message = fmt.Sprintf("File deleted successfully: %s", req.Source)
			}
		case TaskExists:
			exists := fileExists(req.Source)
			log.Printf("File existence check: %s -> %t", req.Source, exists)
			status = "success"
			message = fmt.Sprintf("File %s exists: %t", req.Source, exists)
		default:
			log.Printf("Unknown task type: %s", req.Type)
			status = "fail"
			message = fmt.Sprintf("Unknown task type: %s, request: %+v", req.Type, req)
		}

		// Publish response back to RabbitMQ
		msg := TaskResponse{
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

func deleteFile(src string) error {
	if err := os.Remove(src); err != nil {
		return err
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
