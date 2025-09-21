package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURI      string
	MinioEndpoint    string
	MinioAccessKey   string
	MinioSecretKey   string
	MinieUseSSL      bool
	DBServiceAddress string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		RabbitMQURI:      getEnv("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/"),
		MinioEndpoint:    getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey:   getEnv("MINIO_ACCESS_KEY", "analyze_svc"),
		MinioSecretKey:   getEnv("MINIO_SECRET_KEY", "ANALYZESECRET"),
		MinieUseSSL:      getenvBool("MINIO_USE_SSL"),
		DBServiceAddress: getEnv("DB_SERVICE_ADDRESS", "db-service:50051"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getenvBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
