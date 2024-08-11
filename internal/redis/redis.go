package redisConn

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return rdb
}
