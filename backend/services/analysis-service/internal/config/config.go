package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	RabbitMQURI    string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		RabbitMQURI:    getEnv("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "analyze_svc"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "ANALYZESECRET"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
