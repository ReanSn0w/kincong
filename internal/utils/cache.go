package utils

import (
	"context"
	"sync"
	"time"

	"github.com/go-pkgz/lgr"
)

func NewCache[T any](ctx context.Context, duration time.Duration) *Cache[T] {
	cache := &Cache[T]{
		data: make(map[string]cacheValue[T]),
		mx:   sync.RWMutex{},
	}

	go cache.remover(ctx, duration)
	return cache
}

type cacheValue[T any] struct {
	timestamp time.Time
	value     *T
}

type Cache[T any] struct {
	data map[string]cacheValue[T]
	mx   sync.RWMutex
}

// Get - извлекает значение из хранилища, при доступности
func (c *Cache[T]) Get(key string) (*T, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.data[key]
	return item.value, ok
}

// Must - возвращает значение из кеша или получает
// с помощью замыкания, а затем при необходимости
// добавляет значение в хранилище
func (c *Cache[T]) Must(key string, fn func() (*T, error)) (*T, error) {
	value, ok := c.Get(key)
	if ok {
		return value, nil
	}

	value, err := fn()
	if err != nil {
		return nil, err
	}

	c.Set(key, value)
	return value, nil
}

// Set - сохраняет значение в хранилище
func (c *Cache[T]) Set(key string, value *T) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.data[key] = cacheValue[T]{
		timestamp: time.Now(),
		value:     value}
}

func (c *Cache[T]) remover(ctx context.Context, dur time.Duration) {
	lgr.Default().Logf("[INFO] cache cleaner started")
	ticker := time.NewTicker(time.Second * 10)
	doneCH := ctx.Done()

	for {
		select {
		case <-doneCH:
			lgr.Default().Logf("[INFO] cache cleaner stopped")
			return
		case <-ticker.C:
			lgr.Default().Logf("[DEBUG] cache clean cycle")
			c.removeOldItems(dur)
		}
	}
}

func (c *Cache[T]) removeOldItems(dur time.Duration) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for key, value := range c.data {
		if time.Now().After(value.timestamp.Add(dur)) {
			delete(c.data, key)
		}
	}
}
