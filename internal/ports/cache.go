package ports

import (
	"context"
	"time"
)

type CachePort interface {
	Get(ctx context.Context, cacheKey string) (string, error)
	Set(ctx context.Context, cacheKey string, json []byte, cacheTime time.Duration)
}
