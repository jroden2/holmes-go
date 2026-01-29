package services

import (
	"time"

	"github.com/jroden2/sonic"
)

type cacheService struct {
	sonic sonic.SonicCache
}

func NewCacheService() CacheService {
	return &cacheService{
		sonic: sonic.NewSonicCache(sonic.SonicOptions{
			Capacity: 5,
			TTL:      5 * time.Minute,
		}),
	}
}

type CacheService interface{}

func (c *cacheService) Add(key string, value []byte) {
	c.sonic.Add(key, value)
}

func (c *cacheService) Get(key string) ([]byte, bool) {
	retVal, ok := c.sonic.Get(key)
	if !ok {
		return nil, false
	}

	bytes, ok := retVal.([]byte)
	if !ok {
		return nil, false
	}
	return bytes, true
}

func (c *cacheService) Exists(key string) bool {
	return c.sonic.Exists(key)
}

func (c *cacheService) Purge() {
	c.sonic.Purge()
}

func (c *cacheService) PurgeExpired() {
	c.sonic.PurgeExpired()
}

func (c *cacheService) PeekAll() map[any]any {
	return c.sonic.PeekAll()
}

func (c *cacheService) Close() {
	c.sonic.Close()
}
