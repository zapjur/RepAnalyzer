package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	GRPCPort string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
