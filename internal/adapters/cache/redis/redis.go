package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

type Adapter struct {
	redisClient *redis.Client
}

func NewAdapter(redisAddress, redisPassword string) (*Adapter, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
	})
	return &Adapter{
		redisClient: rdb,
	}, nil
}

func (a Adapter) Get(ctx context.Context, cacheKey string) (string, error) {
	res, err := a.redisClient.Get(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) {
		log.Infof("Cache miss")
		return "", nil
	} else if err != nil {
		return "", errors.New(fmt.Sprintf("error getting from redis: %v\n", err))
	}

	log.Infof("Cache hit")

	return res, nil
}

func (a Adapter) Set(ctx context.Context, cacheKey string, json []byte, cacheTime time.Duration) {
	a.redisClient.Set(ctx, cacheKey, json, cacheTime)
}
