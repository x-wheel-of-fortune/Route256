package in_memory_cache

import (
	"sync"

	"github.com/pkg/errors"

	"homework/internal/pkg/repository"
)

const CacheCapacity = 100

type InMemoryCache struct {
	PickupPoints map[int64]*Node
	mx           sync.RWMutex
	LRU          *Lru
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		PickupPoints: make(map[int64]*Node, CacheCapacity),
		mx:           sync.RWMutex{},
		LRU:          NewLRU(),
	}
	return cache
}

func (c *InMemoryCache) SetPickupPoints(id int64, pickupPoint repository.PickupPoint) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	node_ptr, ok := c.PickupPoints[id]
	if ok {
		c.LRU.set_overwrite(node_ptr, pickupPoint)
		return nil
	}

	if len(c.PickupPoints) == CacheCapacity {
		evictedKey := c.LRU.evict()
		delete(c.PickupPoints, evictedKey)
	}

	nd := &Node{key: id, value: pickupPoint}
	c.PickupPoints[id] = nd

	return nil
}

func (c *InMemoryCache) GetPickupPoints(id int64) (repository.PickupPoint, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	node_ptr, ok := c.PickupPoints[id]
	if !ok {
		return repository.PickupPoint{}, errors.New("cant find pickup point by id")
	}

	c.LRU.get(node_ptr)
	return (*node_ptr).value, nil
}
