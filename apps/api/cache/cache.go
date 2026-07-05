package cache

import (
	"sync"
	"time"
)

// Cache is a simple generic in-memory cache with a max age.
type Cache[T any] struct {
	mu        sync.RWMutex
	data      T
	updatedAt time.Time
	hasData   bool
}

func (c *Cache[T]) Get(maxAge time.Duration) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.hasData || time.Since(c.updatedAt) > maxAge {
		var zero T
		return zero, false
	}
	return c.data, true
}

func (c *Cache[T]) Set(value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = value
	c.updatedAt = time.Now()
	c.hasData = true
}
