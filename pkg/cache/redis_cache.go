package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

// Cache holds a redis connection to the server and connection.
type RedisCache struct {
	c   *cache.Cache
	ttl time.Duration
	log *slog.Logger
}

// NewCache creates a new Cache instance for pgx user to be established.
func NewCache(a, p string, log *slog.Logger) (*RedisCache, error) {
	d, err := getRedisClient(a, p, log)
	if err != nil {
		return nil, err
	}
	return &RedisCache{
		ttl: 10 * time.Second,
		c: cache.New(&cache.Options{
			Redis: d,
		}),
	}, nil
}

func getRedisClient(a, p string, log *slog.Logger) (*redis.Client, error) {
	dsn := fmt.Sprintf("%s:%s", a, p)
	client := redis.NewClient(&redis.Options{
		Addr:               dsn,
		DB:                 1,
		MaxRetries:         5,
		MinRetryBackoff:    5,
		MaxRetryBackoff:    10,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        10 * time.Second,
		PoolTimeout:        10 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 10 * time.Second,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "CacheStore.Connect")
	}
	log.Info("Redis connected")

	return client, nil
}

// Set sets the cache to the given value and returns the corresponding error.
func (cs *RedisCache) Set(ctx context.Context, key string, value any) error {
	d := &cache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          value,
		TTL:            cs.ttl,
		SkipLocalCache: false,
	}
	if err := cs.c.Set(d); err != nil {
		return errors.Wrap(err, "CacheStore.Set")
	}
	return nil
}

// Set sets the cache to the given value and returns the corresponding error.
func (cs *RedisCache) SetWithTimeout(ctx context.Context, key string, value any, t time.Duration) error {
	d := &cache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          value,
		TTL:            t,
		SkipLocalCache: false,
	}
	if err := cs.c.Set(d); err != nil {
		return errors.Wrap(err, "CacheStore.SetWithTimeout")
	}
	return nil
}

// Get returns the cache key for the given key.
func (cs *RedisCache) Get(key string) (any, error) {
	var v any
	if err := cs.c.Get(context.Background(), key, &v); err != nil {
		cs.log.Info(err.Error())
		return v, errors.Wrap(err, "CacheStore.Get")
	}
	return v, nil
}

// Delete removes the specified key from the cache.
func (cs *RedisCache) Delete(key string) error {
	if err := cs.c.Delete(context.Background(), key); err != nil {
		return errors.Wrap(err, "CacheStore.Delete")
	}
	return nil
}
