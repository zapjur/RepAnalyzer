package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisManager struct {
	RedisClient *redis.Client
}

type TaskStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
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

func (rm *RedisManager) SetTaskStatus(videoID, taskName, status string) error {
	ctx := context.Background()
	key := fmt.Sprintf("video:%s:status", videoID)

	payload := TaskStatus{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task status: %w", err)
	}

	return rm.RedisClient.HSet(ctx, key, taskName, jsonData).Err()
}

func (rm *RedisManager) GetTaskStatus(videoID, taskName string) (*TaskStatus, error) {
	ctx := context.Background()
	key := fmt.Sprintf("video:%s:status", videoID)

	val, err := rm.RedisClient.HGet(ctx, key, taskName).Result()
	if err != nil {
		return nil, err
	}

	var status TaskStatus
	if err = json.Unmarshal([]byte(val), &status); err != nil {
		return nil, fmt.Errorf("invalid JSON for task %s: %w", taskName, err)
	}

	return &status, nil
}
