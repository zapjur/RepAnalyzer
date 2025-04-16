package main

import (
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/rabbitmq"
	"orchestrator/internal/redis"
	"orchestrator/internal/server"
)

func main() {

	cfg := config.Load()

	rabbitConn, rabbitChannel, err := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	queues := []string{"bar_path_queue", "bar_path_results_queue"}
	err = rabbitmq.DeclareQueues(rabbitChannel, queues)
	if err != nil {
		log.Fatalf("Failed to declare RabbitMQ queues: %v", err)
	}

	redisClient, err := redis.ConnectToRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	redisManager := &redis.RedisManager{RedisClient: redisClient}

	server.StartGRPCServer(cfg, redisManager, rabbitChannel)
}
