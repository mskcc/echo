package config

import (
	"io"
	"log"
	"os"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	RabbitMQURL       string
	APIToken          string
	TaskQueue         string
	ConfirmationQueue string
	NumberOfWorkers   int
	Port              string
	LogFilePath       string
	LogMaxSize        int
	LogMaxAge         int
	LogMaxBackups     int
}

func Load() (*Config, error) {
	return &Config{
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		APIToken:          getEnv("API_TOKEN", "your-secure-api-token"),
		TaskQueue:         getEnv("FILE_TASK_QUEUE", "file_task_queue"),
		ConfirmationQueue: getEnv("CONFIRMATION_QUEUE", "file_copy_confirmation_queue"),
		NumberOfWorkers:   getEnvInt("NUMBER_OF_WORKERS", 10),
		Port:              getEnv("SERVER_PORT", "8080"),
		LogFilePath:       getEnv("LOG_FILE_PATH", ""),
		LogMaxSize:        getEnvInt("LOG_MAX_SIZE", 1024),
		LogMaxAge:         getEnvInt("LOG_MAX_AGE", 7),
		LogMaxBackups:     getEnvInt("LOG_MAX_BACKUPS", 3),
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

func (c *Config) SetupLogging() error {
	if c.LogFilePath != "" {
		logger := &lumberjack.Logger{
			Filename:   c.LogFilePath,
			MaxSize:    c.LogMaxSize,
			MaxAge:     c.LogMaxAge,
			MaxBackups: c.LogMaxBackups,
			LocalTime:  true,
			Compress:   true,
		}

		multiWriter := io.MultiWriter(os.Stdout, logger)
		log.SetOutput(multiWriter)
	}
	return nil
}
