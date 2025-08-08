package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL       string
	APIToken          string
	TaskQueue         string
	ConfirmationQueue string
	NumberOfWorkers   int
}

func Load() (*Config, error) {
	return &Config{
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		APIToken:          getEnv("API_TOKEN", "your-secure-api-token"),
		TaskQueue:         getEnv("FILE_TASK_QUEUE", "file_task_queue"),
		ConfirmationQueue: getEnv("CONFIRMATION_QUEUE", "file_copy_confirmation_queue"),
		NumberOfWorkers:   getEnvInt("NUMBER_OF_WORKERS", 10),
	}, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}
