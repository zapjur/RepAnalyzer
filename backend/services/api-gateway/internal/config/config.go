package config

import "os"

type Config struct {
	HTTPPort           string
	UserServiceAddress string
}

func Load() *Config {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "user-service:50051"
	}

	return &Config{
		HTTPPort:           port,
		UserServiceAddress: userServiceAddr,
	}
}
