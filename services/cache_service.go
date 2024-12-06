package services

import (
	"encoding/json"
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

func (i CacheItem) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

type CacheService struct {
	mu             sync.RWMutex
	items          map[string]CacheItem
	cleanupTicker  *time.Ticker
	cleanupRunning bool
	defaultTTL     time.Duration
}

func NewCacheService(defaultTTL time.Duration, cleanupInterval time.Duration) *CacheService {
	cache := &CacheService{
		items:      make(map[string]CacheItem),
		defaultTTL: defaultTTL,
	}

	if cleanupInterval > 0 {
		cache.cleanupTicker = time.NewTicker(cleanupInterval)
		go cache.startCleanupTimer()
	}

	return cache
}

func (c *CacheService) startCleanupTimer() {
	c.mu.Lock()
	if c.cleanupRunning {
		c.mu.Unlock()
		return
	}
	c.cleanupRunning = true
	c.mu.Unlock()

	for range c.cleanupTicker.C {
		c.DeleteExpired()
	}
}

func (c *CacheService) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

func (c *CacheService) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: exp,
	}
	c.mu.Unlock()
}

func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	if item.Expired() {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	return item.Value, true
}

func GetTyped[T any](c *CacheService, key string) (T, bool) {
	var result T
	value, found := c.Get(key)
	if !found {
		return result, false
	}

	if typed, ok := value.(T); ok {
		return typed, true
	}

	data, err := json.Marshal(value)
	if err != nil {
		return result, false
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, false
	}

	return result, true
}

func (c *CacheService) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *CacheService) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

func (c *CacheService) Clear() {
	c.mu.Lock()
	c.items = make(map[string]CacheItem)
	c.mu.Unlock()
}

func (c *CacheService) Count() int {
	c.mu.RLock()
	count := len(c.items)
	c.mu.RUnlock()
	return count
}

func (c *CacheService) Keys() []string {
	c.mu.RLock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	c.mu.RUnlock()
	return keys
}

func (c *CacheService) SetDefaultTTL(ttl time.Duration) {
	c.defaultTTL = ttl
}

func (c *CacheService) GetDefaultTTL() time.Duration {
	return c.defaultTTL
}

func (c *CacheService) Close() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
}
