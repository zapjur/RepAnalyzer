package redis

import (
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	redisClient *redis.Client
}
