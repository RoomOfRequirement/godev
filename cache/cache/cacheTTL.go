package cache

import (
	"fmt"
	"godev/cache"
	"godev/cache/lru"
	"os"
	"sync"
	"time"
)

// TTL struct
type TTL struct {
	lru      *lru.LRU
	lock     sync.RWMutex
	ttl      time.Duration
	stopChan chan bool
}

// ttlEntry stored inside lru item value
type ttlEntry struct {
	value            interface{}
	lastAccessedTime time.Time
}

// NewCacheTTL creates a new ttl-enabled lru cache with input size
//	notice: default least ttl is 1 second, default least cleanInterval is 2 seconds
func NewCacheTTL(size int, ttl, cleanInterval time.Duration) (*TTL, error) {
	return NewCacheTTLWithOnEvict(size, ttl, cleanInterval, nil)
}

// NewCacheTTLWithOnEvict creates a new ttl-enabled lru cache with input size and onEvict function
//	notice: default least ttl is 1 second, default least cleanInterval is 2 seconds
func NewCacheTTLWithOnEvict(size int, ttl, cleanInterval time.Duration, onEvict cache.EvictCallback) (*TTL, error) {
	c, err := lru.NewLRU(size, onEvict)
	if err != nil {
		return nil, err
	}

	if ttl <= 1*time.Second {
		ttl = 1 * time.Second
	}

	cacheTTL := &TTL{
		lru:  c,
		lock: sync.RWMutex{},
		ttl:  ttl,
	}

	if cleanInterval <= 1*time.Second {
		cleanInterval = 2 * time.Second
	}

	cacheTTL.stopChan = cacheTTL.removeExpired(time.NewTicker(cleanInterval))
	return cacheTTL, nil
}

// Add adds key value to the cache if not found else update value and returns true if an eviction occurred
func (c *TTL) Add(key, value interface{}) (found, evicted bool) {
	c.lock.Lock()
	found, evicted = c.lru.Add(key, ttlEntry{
		value:            value,
		lastAccessedTime: time.Now(),
	})
	c.lock.Unlock()
	return
}

// Get returns value if key and update timestamp
func (c *TTL) Get(key interface{}) (value interface{}, found bool) {
	// Lock due to updating timestamp
	c.lock.Lock()
	entry, found := c.lru.Get(key)
	if found {
		// lastAccessedTime + ttl > time.now() => evicted => remove key value pair and set found to false
		if lastAccessedTime := entry.(ttlEntry).lastAccessedTime; time.Now().After(lastAccessedTime.Add(c.ttl)) {
			c.lru.Remove(key)
			found = false
		} else {
			value = entry.(ttlEntry).value
		}
	}
	c.lock.Unlock()
	return
}

// Peek returns value of key without updating timestamp
func (c *TTL) Peek(key interface{}) (value interface{}, found bool) {
	// RLock due to not updating
	c.lock.RLock()
	entry, found := c.lru.Peek(key)
	if found {
		value = entry.(ttlEntry).value
	}
	c.lock.RUnlock()
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (c *TTL) Contains(key interface{}) (found bool) {
	c.lock.RLock()
	found = c.lru.Contains(key)
	c.lock.RUnlock()
	return
}

// Remove removes key value if found in the cache
func (c *TTL) Remove(key interface{}) (value interface{}, found bool) {
	c.lock.Lock()
	entry, found := c.lru.Remove(key)
	if found {
		value = entry.(ttlEntry).value
	}
	c.lock.Unlock()
	return
}

// GetLeastUsed returns least used key value pairs if found in the cache
func (c *TTL) GetLeastUsed() (key, value interface{}, found bool) {
	c.lock.RLock()
	key, entry, found := c.lru.GetLeastUsed()
	if found {
		value = entry.(ttlEntry).value
	}
	c.lock.RUnlock()
	return
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
func (c *TTL) RemoveLeastUsed() (key, value interface{}, found bool) {
	c.lock.Lock()
	key, entry, found := c.lru.RemoveLeastUsed()
	if found {
		value = entry.(ttlEntry).value
	}
	c.lock.Unlock()
	return
}

// Size returns key value pair numbers in the cache
//	if want memory size in bytes, need to use unsafe.Sizeof, which may not be supported in some platform
func (c *TTL) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Size()
}

// Resize resize cache size and returns diff
func (c *TTL) Resize(size int) (diff int, err error) {
	c.lock.Lock()
	diff, err = c.lru.Resize(size)
	c.lock.Unlock()
	return
}

// Clear clears all pairs in the cache
func (c *TTL) Clear() {
	c.lock.Lock()
	c.lru.Clear()
	c.lock.Unlock()
	return
}

// ResetTTL resets ttl
func (c *TTL) ResetTTL(ttl time.Duration) {
	c.lock.Lock()
	c.ttl = ttl
	c.lock.Unlock()
	return
}

// StopCleanWork stops clean expired goroutine
func (c *TTL) StopCleanWork() {
	c.lock.Lock()
	c.stopChan <- true
	close(c.stopChan)
	c.lock.Unlock()
	return
}

// RestartCleanWork restarts clean expired goroutine
func (c *TTL) RestartCleanWork(cleanInterval time.Duration) {
	if cleanInterval <= 1*time.Second {
		cleanInterval = 2 * time.Second
	}
	// !ok => closed
	c.lock.Lock()
	if _, ok := <-c.stopChan; !ok {
		c.removeExpired(time.NewTicker(cleanInterval))
	} else {
		c.stopChan <- true
		close(c.stopChan)
		c.removeExpired(time.NewTicker(cleanInterval))
	}
	c.lock.Unlock()
	return
}

// removeExpired removes expired key pairs in one background goroutine
//	notice: default cleanInterval is 2 seconds
//	https://stackoverflow.com/questions/17797754/ticker-stop-behaviour-in-golang
func (c *TTL) removeExpired(cleanInterval *time.Ticker) chan bool {
	stopChan := make(chan bool, 1)

	go func(cleanInterval *time.Ticker) {
		defer cleanInterval.Stop()

		for {
			select {
			case <-cleanInterval.C:
				for _, key := range c.lru.Keys() {
					c.lock.Lock()
					entry, found := c.lru.Get(key)
					c.lock.Unlock()
					if found && time.Now().After(entry.(ttlEntry).lastAccessedTime.Add(c.ttl)) {
						c.Remove(key)
					}
				}
			case stop := <-stopChan:
				if stop {
					_, _ = fmt.Fprintln(os.Stderr, "clean expired goroutine stopped")
					return
				}
			}
		}
	}(cleanInterval)
	return stopChan
}
