package list

import (
	container "goContainer"
	"math/rand"
)

type node struct {
	key, value interface{}
	next       []*node // next levels
}

// SkipList struct
type SkipList struct {
	maxLevel   int
	decFactor  int // should be power of 2
	header     *node
	comparator container.Comparator
	itemNum    int
}

// NewSkipList creates a new skip list according to input arguments
//	notice: skip level starts from 0, level 0 is a common linked list
//	decFactor is used to determine the possibility new node whether should have next level pointer
//	comparator for comparision of key
func NewSkipList(maxLevel, decFactor int, comparator container.Comparator) *SkipList {
	if maxLevel < 0 {
		maxLevel = 0
	}
	if decFactor < 2 {
		decFactor = 2
	} else {
		decFactor = nextPowerOfTwo(decFactor)
	}
	return &SkipList{
		maxLevel:  maxLevel,  // 2 ^ maxLevel items
		decFactor: decFactor, // decrease factor for every layer from bottom level to top level
		header: &node{
			key:   nil,
			value: nil,
			next:  []*node{nil},
		}, // top level first node
		comparator: comparator,
		itemNum:    1, // header
	}
}

func nextPowerOfTwo(n int) int {
	if n > 0 && n&(n-1) == 0 {
		return n
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16 // 32 bit OS, runtime.GOARCH
	n |= n >> 32 // 64 bit OS
	n++
	return n
}

func (sl *SkipList) randomLevel() int {
	// rand.Seed(time.Now().UnixNano())
	k := 0 // level 0 is the bottom linked list
	for rand.Int()&sl.decFactor != 0 && k < sl.maxLevel {
		k++
	}
	return k
}

// Set insert k, v into skip list or update v of k
func (sl *SkipList) Set(k, v interface{}) {
	foundNode, path := sl.searchPath(sl.header, k)
	// if found key already inside the skip list, just update its value
	if foundNode != nil && sl.comparator(foundNode.key, k) == 0 {
		foundNode.value = v
		return
	}

	newLevel := sl.randomLevel()
	currentLevel := len(sl.header.next) - 1
	if newLevel > currentLevel {
		for i := currentLevel + 1; i <= newLevel; i++ {
			// record new node path
			path = append(path, sl.header)
			// update header, header should be in all levels
			sl.header.next = append(sl.header.next, nil)
		}
	}

	nNode := &node{
		key:   k,
		value: v,
		next:  make([]*node, newLevel+1, sl.maxLevel+1),
	}

	for i := 0; i <= newLevel; i++ {
		nNode.next[i] = path[i].next[i]
		path[i].next[i] = nNode
	}

	sl.itemNum++
}

// searchPath returns node with the same key or the upper bound node
func (sl *SkipList) searchPath(current *node, key interface{}) (*node, []*node) {
	path := make([]*node, len(sl.header.next), sl.maxLevel+1)
	// search from current level to the bottom (level 0)
	for i := len(current.next) - 1; i >= 0; i-- {
		// find the proper level where key <= current.key
		for current.next[i] != nil && sl.comparator(current.next[i].key, key) < 0 {
			current = current.next[i]
		}
		path[i] = current
	}
	if len(current.next) == 0 {
		return nil, path
	}
	return current.next[0], path
}

// Search returns true if k found in skip list
func (sl *SkipList) Search(k interface{}) bool {
	node, _ := sl.searchPath(sl.header, k)

	if node == nil || node.key != k {
		return false
	}

	return true
}

// Get returns v with input k if found k in skip list
func (sl *SkipList) Get(k interface{}) (v interface{}, found bool) {
	node, _ := sl.searchPath(sl.header, k)

	if node == nil || node.key != k {
		return nil, false
	}

	return node.value, true
}

// Delete deletes node with input k and return its v
//	`ok` indicates whether deletion succeeds or not
func (sl *SkipList) Delete(k interface{}) (v interface{}, ok bool) {
	node, path := sl.searchPath(sl.header, k)
	// k not inside the skip list
	if node == nil || node.key == nil {
		return nil, false
	}

	for i := 0; i <= len(sl.header.next)-1 && path[i].next[i] == node; i++ {
		path[i].next[i] = node.next[i]
	}

	for len(sl.header.next)-1 > 0 && sl.header.next[len(sl.header.next)-1] == nil {
		sl.header.next = sl.header.next[:len(sl.header.next)-1]
	}
	sl.itemNum--

	return node.value, true
}

// Empty returns true if no k, v stored inside skip list
func (sl *SkipList) Empty() bool {
	return sl.itemNum-1 == 0
}

// Size returns the quantity of k, v stored inside skip list
func (sl *SkipList) Size() int {
	return sl.itemNum - 1
}

// Clear clears skip list
func (sl *SkipList) Clear() {
	sl.header = nil
	sl.itemNum = 1
}

// Values returns k, v inside skip list
func (sl *SkipList) Values() []interface{} {
	current := sl.header
	values := make([]interface{}, 0, sl.itemNum)
	for i := 0; i < sl.itemNum; i++ {
		if current.next == nil || len(current.next) == 0 {
			break
		}
		values = append(values, current.value)
		current = current.next[0]
	}
	// skip nil header
	return values[1:]
}
