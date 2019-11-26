package cache

// Interface of cache
//	https://en.wikipedia.org/wiki/Cache_replacement_policies
type Interface interface {
	// Add adds key value to the cache if not found else update value and returns true if an eviction occurred
	Add(key, value interface{}) (found, evicted bool)
	// Get returns value if key and update timestamp
	Get(key interface{}) (value interface{}, found bool)
	// Peek returns value of key without updating timestamp
	Peek(key interface{}) (value interface{}, found bool)
	// Contains returns true if key found in the cache without updating timestamp
	Contains(key interface{}) (found bool)
	// Remove removes key value if found in the cache
	Remove(key interface{}) (value interface{}, found bool)
	// Size returns key value pair numbers in the cache
	//	if want memory size in bytes, need to use unsafe.Sizeof, which may not be supported in some platform
	Size() int
	// Resize resize cache size and returns diff
	Resize(size int) (diff int, err error)

	// Clear clears all pairs in the cache
	Clear()
}
