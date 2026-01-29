package services

import (
	"time"

	"github.com/jroden2/sonic"
)

type KVP sonic.Entry

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

func (c *cacheService) Add(kvp KVP) {

}
