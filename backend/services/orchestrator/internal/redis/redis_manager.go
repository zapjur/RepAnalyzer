package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisClient *redis.Client
}

func ConnectToRedis(addr string) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Connected to Redis successfully.")
	return client, nil
}
