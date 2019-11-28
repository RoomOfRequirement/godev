package lfu

import (
	"container/list"
	"errors"
	"goContainer/cache"
)

// LFU struct (least frequently used)
//	http://dhruvbird.com/lfu.pdf
type LFU struct {
	size int

	freqList *list.List
	pairs    map[interface{}]*item

	onEvict cache.EvictCallback
}

type freqItem struct {
	freq int
	// instead of using doubly linked list in the paper,
	// choose map since all items in this freqItem have the same freq
	items map[*item]struct{}
}

type item struct {
	freqListPos *list.Element
	key         interface{}
	value       interface{}
}

// NewLFU creates a new LFU cache with input size and onEvict function
func NewLFU(size int, onEvict cache.EvictCallback) (*LFU, error) {
	if size <= 0 {
		return nil, errors.New("invalid cache size")
	}
	return &LFU{
		size:     size,
		freqList: list.New(),
		pairs:    make(map[interface{}]*item),
		onEvict:  onEvict,
	}, nil
}

// Add adds key value to the cache and returns true if an eviction occurred
func (lfu *LFU) Add(key, value interface{}) (found, evicted bool) {
	// if found, update value
	if it, found := lfu.pairs[key]; found {
		lfu.updateFreq(it)
		it.value = value
		return true, false
	}

	// not found
	// check whether need evict least used
	evicted = len(lfu.pairs) >= lfu.size
	// eviction occurred
	if evicted {
		lfu.removeLeastUsed()
	}

	// add new item
	nit := &item{
		freqListPos: nil,
		key:         key,
		value:       value,
	}
	lfu.updateFreq(nit)
	lfu.pairs[key] = nit

	return false, evicted
}

// Get returns value if key and update timestamp
func (lfu *LFU) Get(key interface{}) (value interface{}, found bool) {
	it, found := lfu.pairs[key]
	if !found {
		return
	}
	value = it.value
	lfu.updateFreq(it)
	return
}

// Peek returns value of key without updating timestamp
func (lfu *LFU) Peek(key interface{}) (value interface{}, found bool) {
	it, found := lfu.pairs[key]
	if !found {
		return
	}
	value = it.value
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (lfu *LFU) Contains(key interface{}) (found bool) {
	_, found = lfu.pairs[key]
	return
}

// Remove removes key value if found in the cache
func (lfu *LFU) Remove(key interface{}) (value interface{}, found bool) {
	it, found := lfu.pairs[key]
	if found {
		lfu.removeItem(it.freqListPos, it)
		value = it.value
		delete(lfu.pairs, it.key)
		return
	}
	return
}

// Size returns key value pair numbers in the cache
func (lfu *LFU) Size() int {
	return len(lfu.pairs)
}

// Resize resize cache size and returns diff
func (lfu *LFU) Resize(size int) (diff int, err error) {
	diff = size - lfu.size
	if size > 0 {
		lfu.size = size
	} else {
		return 0, errors.New("resize size should be larger than 0")
	}

	if diff >= 0 {
		return
	}

	for i := 0; i < -diff; i++ {
		lfu.removeLeastUsed()
	}
	return
}

// Clear clears all pairs in the cache
func (lfu *LFU) Clear() {
	for k, it := range lfu.pairs {
		if lfu.onEvict != nil {
			lfu.onEvict(k, it.value)
		}
		delete(lfu.pairs, k)
	}
	lfu.freqList.Init()
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
func (lfu *LFU) RemoveLeastUsed() (key, value interface{}, found bool) {
	fi := lfu.freqList.Front()
	if fi == nil {
		return
	}
	for it := range fi.Value.(*freqItem).items {
		key, value = it.key, it.value
		found = true
		if lfu.onEvict != nil {
			lfu.onEvict(it.key, it.value)
		}
		delete(lfu.pairs, it.key)
		lfu.removeItem(fi, it)
		break
	}
	return
}

// since map has no order, it will remove random one with least freq
func (lfu *LFU) removeLeastUsed() {
	fi := lfu.freqList.Front()
	if fi == nil {
		return
	}
	for it := range fi.Value.(*freqItem).items {
		if lfu.onEvict != nil {
			lfu.onEvict(it.key, it.value)
		}
		delete(lfu.pairs, it.key)
		lfu.removeItem(fi, it)
		break
	}
}

func (lfu *LFU) updateFreq(it *item) {
	currentPos := it.freqListPos
	var nextFreq int
	var nextPos *list.Element

	// new item
	if currentPos == nil {
		nextFreq = 1
		nextPos = lfu.freqList.Front()
	} else {
		nextFreq = currentPos.Value.(*freqItem).freq + 1
		nextPos = currentPos.Next()
	}

	// next freq list pos not exist or its freq number does not meet requirement
	// create new freq pos as next
	if nextPos == nil || nextPos.Value.(*freqItem).freq != nextFreq {
		newFreqItem := &freqItem{
			freq:  nextFreq,
			items: make(map[*item]struct{}),
		}

		// new item
		if currentPos == nil {
			nextPos = lfu.freqList.PushFront(newFreqItem)
		} else {
			nextPos = lfu.freqList.InsertAfter(newFreqItem, currentPos)
		}
	}

	// update freqListPos
	it.freqListPos = nextPos
	// update nextPos
	nextPos.Value.(*freqItem).items[it] = struct{}{}
	// update currentPos
	if currentPos != nil {
		lfu.removeItem(currentPos, it)
	}
}

func (lfu *LFU) removeItem(pos *list.Element, it *item) {
	items := pos.Value.(*freqItem).items
	delete(items, it)
	if len(items) == 0 {
		lfu.freqList.Remove(pos)
	}
	return
}
