package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	GRPCPort         string
	RabbitMQURI      string
	RedisAddr        string
	DBServiceAddress string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		GRPCPort:         getEnv("GRPC_PORT", "50051"),
		RabbitMQURI:      getEnv("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/"),
		RedisAddr:        getEnv("REDIS_ADDR", "redis:6379"),
		DBServiceAddress: getEnv("DB_SERVICE_ADDRESS", "db-service:50051"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
