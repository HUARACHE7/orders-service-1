package cache

import (
	"sync"
	"github.com/HUARACHE7/orders-service-1/internal/model"
)

type AppCache struct {
	mu    sync.RWMutex 
	items map[string]model.Order
}

func NewCache() *AppCache {
	return &AppCache{
		items: make(map[string]model.Order),
	}
}

func (c *AppCache) Set(key string, order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = order
}

func (c *AppCache) Get(key string) (model.Order, bool) {
	c.mu.RLock() 
	defer c.mu.RUnlock()
	item, found := c.items[key]
	return item, found
}

func (c *AppCache) Load(items map[string]model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = items
}
