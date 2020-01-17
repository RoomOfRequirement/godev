package cache

import (
	"errors"
	"godev/cache"
	"godev/cache/arc"
	"godev/cache/lfu"
	"godev/cache/lru"
	"sync"
)

// Cache struct
type Cache struct {
	cache cache.Interface
	lock  sync.RWMutex
}

// Type for backend cache type
type Type int

// supported backend cache type
const (
	LRU = iota
	LFU
	ARC
)

// NewCache creates a new lru cache with input size
func NewCache(size int, cacheType Type) (*Cache, error) {
	return NewCacheWithOnEvict(size, cacheType, nil)
}

// NewCacheWithOnEvict creates a new lru cache with input size and onEvict function
func NewCacheWithOnEvict(size int, cacheType Type, onEvict cache.EvictCallback) (*Cache, error) {
	var err error
	var c cache.Interface
	switch cacheType {
	case LRU:
		c, err = lru.NewLRU(size, onEvict)
	case LFU:
		c, err = lfu.NewLFU(size, onEvict)
	case ARC:
		c, err = arc.NewARC(size, onEvict)
	default:
		err = errors.New("unsupported cache type")
	}
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: c,
		lock:  sync.RWMutex{},
	}, nil
}

// Add adds key value to the cache if not found else update value and returns true if an eviction occurred
func (c *Cache) Add(key, value interface{}) (found, evicted bool) {
	c.lock.Lock()
	found, evicted = c.cache.Add(key, value)
	c.lock.Unlock()
	return
}

// Get returns value if key and update timestamp
func (c *Cache) Get(key interface{}) (value interface{}, found bool) {
	// Lock due to updating timestamp
	c.lock.Lock()
	value, found = c.cache.Get(key)
	c.lock.Unlock()
	return
}

// Peek returns value of key without updating timestamp
func (c *Cache) Peek(key interface{}) (value interface{}, found bool) {
	// RLock due to not updating
	c.lock.RLock()
	value, found = c.cache.Peek(key)
	c.lock.RUnlock()
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (c *Cache) Contains(key interface{}) (found bool) {
	c.lock.RLock()
	found = c.cache.Contains(key)
	c.lock.RUnlock()
	return
}

// Remove removes key value if found in the cache
func (c *Cache) Remove(key interface{}) (value interface{}, found bool) {
	c.lock.Lock()
	value, found = c.cache.Remove(key)
	c.lock.Unlock()
	return
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
func (c *Cache) RemoveLeastUsed() (key, value interface{}, found bool) {
	c.lock.Lock()
	key, value, found = c.cache.RemoveLeastUsed()
	c.lock.Unlock()
	return
}

// Size returns key value pair numbers in the cache
//	if want memory size in bytes, need to use unsafe.Sizeof, which may not be supported in some platform
func (c *Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.Size()
}

// Resize resize cache size and returns diff
func (c *Cache) Resize(size int) (diff int, err error) {
	c.lock.Lock()
	diff, err = c.cache.Resize(size)
	c.lock.Unlock()
	return
}

// Clear clears all pairs in the cache
func (c *Cache) Clear() {
	c.lock.Lock()
	c.cache.Clear()
	c.lock.Unlock()
	return
}
