package gomemorycache

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/zhashkevych/scheduler"
)

type GoMemoryCache struct {
	mx         sync.RWMutex
	cacheItems map[string]CacheItems
}

type CacheItems struct {
	Value          interface{}
	ExpirationTime time.Time
}

// constructor
func New() *GoMemoryCache {
	ctx := context.Background()
	cache := GoMemoryCache{
		cacheItems: make(map[string]CacheItems),
	}
	worker := scheduler.NewScheduler()
	worker.Add(ctx, cache.cleanLoop, time.Second*1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
	worker.Stop()

	return &cache
}

func (c *GoMemoryCache) cleanLoop(ctx context.Context) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if len(c.cacheItems) > 0 {
		for key, item := range c.cacheItems {
			if item.ExpirationTime.Unix() < time.Now().Unix() {
				delete(c.cacheItems, key)
			}
		}
	}
}

func (c *GoMemoryCache) Get(key string) (interface{}, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.cacheItems[key]

	if !ok {
		return nil, errors.New("record not found")
	}

	return item.Value, nil
}

func (c *GoMemoryCache) Set(key string, val interface{}, ttl time.Duration) {
	c.mx.Lock()
	expTime := time.Now().Add(ttl)
	c.cacheItems[key] = CacheItems{
		Value:          val,
		ExpirationTime: expTime,
	}
	c.mx.Unlock()
}

func (c *GoMemoryCache) Delete(key string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, exist := c.cacheItems[key]; !exist {
		return errors.New("record not found")
	}

	delete(c.cacheItems, key)

	return nil
}
