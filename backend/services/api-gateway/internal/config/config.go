package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort            string
	Auth0Domain         string
	HTTPPort            string
	DBServiceAddress    string
	OrchestratorAddress string
}

func Load() *Config {
	_ = godotenv.Load()

	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("HTTP_PORT", "8080")
	auth0Domain := getEnv("AUTH0_DOMAIN", "")
	dbServiceAddress := getEnv("DB_SERVICE_ADDRESS", "db-service:50051")
	orchestratorAddress := getEnv("ORCHESTRATOR_SERVICE_ADDRESS", "orchestrator:50051")

	log.Println("AUTH0_DOMAIN:", auth0Domain)

	return &Config{
		GRPCPort:            grpcPort,
		HTTPPort:            httpPort,
		Auth0Domain:         auth0Domain,
		DBServiceAddress:    dbServiceAddress,
		OrchestratorAddress: orchestratorAddress,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
