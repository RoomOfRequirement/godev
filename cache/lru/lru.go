package lru

import (
	"container/list"
	"errors"
	"goContainer/cache"
)

// LRU struct
type LRU struct {
	size int

	// evicts and pairs consist of linked hash map
	//	see https://github.com/Harold2017/GoContainer/blob/master/maps/linkedhashmap/linkedhashmap.go
	evicts *list.List
	pairs  map[interface{}]*item

	onEvict cache.EvictCallback
}

type item struct {
	key   *list.Element
	value interface{}
}

// NewLRU creates a new LRU cache with input size and onEvict function
func NewLRU(size int, onEvict cache.EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("invalid cache size")
	}
	return &LRU{
		size:    size,
		evicts:  list.New(),
		pairs:   make(map[interface{}]*item),
		onEvict: onEvict,
	}, nil
}

// Add adds key value to the cache and returns true if an eviction occurred
func (lru *LRU) Add(key, value interface{}) (found, evicted bool) {
	// if found, update value
	if e, found := lru.pairs[key]; found {
		lru.evicts.MoveToFront(e.key)
		e.value = value
		return true, false
	}

	// not found, add new item
	e := lru.evicts.PushFront(key)
	lru.pairs[key] = &item{
		key:   e,
		value: value,
	}

	evicted = lru.evicts.Len() > lru.size
	// eviction occurred
	if evicted {
		lru.removeLeastUsed()
		return false, true
	}
	return false, false
}

// Get returns value if key and update timestamp
func (lru *LRU) Get(key interface{}) (value interface{}, found bool) {
	it, found := lru.pairs[key]
	if !found {
		return
	}
	value = it.value
	lru.evicts.MoveToFront(it.key)
	return
}

// Peek returns value of key without updating timestamp
func (lru *LRU) Peek(key interface{}) (value interface{}, found bool) {
	it, found := lru.pairs[key]
	if !found {
		return
	}
	value = it.value
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (lru *LRU) Contains(key interface{}) (found bool) {
	_, found = lru.pairs[key]
	return
}

// Remove removes key value if found in the cache
func (lru *LRU) Remove(key interface{}) (value interface{}, found bool) {
	it, found := lru.pairs[key]
	if found {
		lru.evicts.Remove(it.key)
		value = lru.pairs[key].value
		delete(lru.pairs, key)
		return
	}
	return
}

// Size returns key value pair numbers in the cache
//	if want memory size in bytes, need to use unsafe.Sizeof, which may not be supported in some platform
func (lru *LRU) Size() int {
	return lru.evicts.Len()
}

// Resize resize cache size and returns diff
func (lru *LRU) Resize(size int) (diff int, err error) {
	diff = size - lru.size
	if size > 0 {
		lru.size = size
	} else {
		return 0, errors.New("resize size should be larger than 0")
	}

	if diff >= 0 {
		return
	}

	for i := 0; i < -diff; i++ {
		lru.removeLeastUsed()
	}
	return
}

// Clear clears all pairs in the cache
func (lru *LRU) Clear() {
	for k, it := range lru.pairs {
		if lru.onEvict != nil {
			lru.onEvict(k, it.value)
		}
		delete(lru.pairs, k)
	}
	lru.evicts.Init()
}

// GetLeastUsed returns least used key value pairs if found in the cache
func (lru *LRU) GetLeastUsed() (key, value interface{}, found bool) {
	e := lru.evicts.Back()
	if e != nil {
		key = e.Value
		value = lru.pairs[key].value
		found = true
	}
	return
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
func (lru *LRU) RemoveLeastUsed() (key, value interface{}, found bool) {
	e := lru.evicts.Back()
	key, value = lru.removeItem(e)
	if key != nil {
		found = true
	}
	return
}

func (lru *LRU) removeLeastUsed() {
	e := lru.evicts.Back()
	lru.removeItem(e)
}

func (lru *LRU) removeItem(e *list.Element) (key, value interface{}) {
	if e != nil {
		key = e.Value
		value = lru.pairs[key].value
		lru.evicts.Remove(e)
		delete(lru.pairs, key)
		if lru.onEvict != nil {
			lru.onEvict(key, value)
		}
	}
	return
}

// Keys returns all keys in the cache, from oldest to newest
func (lru *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(lru.pairs))
	i := 0
	for e := lru.evicts.Back(); e != nil; e = e.Prev() {
		keys[i] = e.Value
		i++
	}
	return keys
}
