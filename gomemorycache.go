package gomemorycache

import (
	"errors"
	"sync"
)

type GoMemoryCache struct {
	mx         sync.RWMutex
	cacheItems map[string]CacheItems
}

type CacheItems struct {
	Value interface{}
}

// constructor
func New() *GoMemoryCache {
	return &GoMemoryCache{
		cacheItems: make(map[string]CacheItems),
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

func (c *GoMemoryCache) Set(key string, val interface{}) {
	c.mx.Lock()
	c.cacheItems[key] = CacheItems{Value: val}
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
