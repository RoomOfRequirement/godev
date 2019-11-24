package lru

import "sync"

// Cache struct
type Cache struct {
	lru  *LRU
	lock sync.RWMutex
}

// NewCache creates a new lru cache with input size
func NewCache(size int) (*Cache, error) {
	return NewCacheWithOnEvict(size, nil)
}

// NewCacheWithOnEvict creates a new lru cache with input size and onEvict function
func NewCacheWithOnEvict(size int, onEvict EvictCallback) (*Cache, error) {
	lru, err := NewLRU(size, onEvict)
	if err != nil {
		return nil, err
	}
	return &Cache{
		lru:  lru,
		lock: sync.RWMutex{},
	}, nil
}

// Add adds key value to the cache if not found else update value and returns true if an eviction occurred
func (c *Cache) Add(key, value interface{}) (found, evicted bool) {
	c.lock.Lock()
	found, evicted = c.lru.Add(key, value)
	c.lock.Unlock()
	return
}

// Get returns value if key and update timestamp
func (c *Cache) Get(key interface{}) (value interface{}, found bool) {
	// Lock due to updating timestamp
	c.lock.Lock()
	value, found = c.lru.Get(key)
	c.lock.Unlock()
	return
}

// Peek returns value of key without updating timestamp
func (c *Cache) Peek(key interface{}) (value interface{}, found bool) {
	// RLock due to not updating
	c.lock.RLock()
	value, found = c.lru.Peek(key)
	c.lock.RUnlock()
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (c *Cache) Contains(key interface{}) (found bool) {
	c.lock.RLock()
	found = c.lru.Contains(key)
	c.lock.RUnlock()
	return
}

// Remove removes key value if found in the cache
func (c *Cache) Remove(key interface{}) (value interface{}, found bool) {
	c.lock.Lock()
	value, found = c.lru.Remove(key)
	c.lock.Unlock()
	return
}

// GetLeastUsed returns least used key value pairs if found in the cache
func (c *Cache) GetLeastUsed() (key, value interface{}, found bool) {
	c.lock.RLock()
	key, value, found = c.lru.GetLeastUsed()
	c.lock.RUnlock()
	return
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
func (c *Cache) RemoveLeastUsed() (key, value interface{}, found bool) {
	c.lock.Lock()
	key, value, found = c.lru.RemoveLeastUsed()
	c.lock.Unlock()
	return
}

// Size returns key value pair numbers in the cache
//	if want memory size in bytes, need to use unsafe.Sizeof, which may not be supported in some platform
func (c *Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Size()
}

// Resize resize cache size and returns diff
func (c *Cache) Resize(size int) (diff int) {
	c.lock.Lock()
	diff = c.lru.Resize(size)
	c.lock.Unlock()
	return
}

// Clear clears all pairs in the cache
func (c *Cache) Clear() {
	c.lock.Lock()
	c.lru.Clear()
	c.lock.Unlock()
	return
}
