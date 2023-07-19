package cache

import (
	"context"
	"sync"
	"time"
)

// LocalCache Cache implements the Store interface.
type LocalCache struct {
	data     map[string]any
	mutex    sync.RWMutex
	expiry   map[string]time.Time
	expMutex sync.RWMutex
}

// NewLocalCache NewCache creates a new Cache instance.
func NewLocalCache() *LocalCache {
	return &LocalCache{
		data:     make(map[string]any),
		expiry:   make(map[string]time.Time),
		mutex:    sync.RWMutex{},
		expMutex: sync.RWMutex{},
	}
}

// Set stores a value in the cache with the given key.
func (c *LocalCache) Set(ctx context.Context, key string, value any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return nil
}

// SetWithTimeout stores a value in the cache with the given key and expiration time.
func (c *LocalCache) SetWithTimeout(ctx context.Context, key string, value any, t time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	c.expMutex.Lock()
	defer c.expMutex.Unlock()
	c.expiry[key] = time.Now().Add(t)
	return nil
}

// Get retrieves a value from the cache based on the given key.
func (c *LocalCache) Get(key string) (any, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, nil // Key does not exist in cache
	}
	c.expMutex.RLock()
	defer c.expMutex.RUnlock()
	if expiry, ok := c.expiry[key]; ok && expiry.Before(time.Now()) {
		delete(c.data, key)   // Delete expired value
		delete(c.expiry, key) // Delete expired expiry time
		return nil, nil       // Key has expired
	}
	return value, nil
}

// Delete removes a value from the cache based on the given key.
func (c *LocalCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	c.expMutex.Lock()
	defer c.expMutex.Unlock()
	delete(c.expiry, key)
	return nil
}
