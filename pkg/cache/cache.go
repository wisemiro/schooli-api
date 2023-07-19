package cache

import (
	"context"
	"time"
)

type Store interface {
	Set(ctx context.Context, key string, value any) error
	SetWithTimeout(ctx context.Context, key string, value any, t time.Duration) error
	Get(key string) (any, error)
	Delete(key string) error
}
