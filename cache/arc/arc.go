package arc

import (
	"container/list"
	"errors"
)

// ARC struct for adaptive replacement cache
//	https://en.wikipedia.org/wiki/Adaptive_replacement_cache
//	http://u.cs.biu.ac.il/~wiseman/2os/2os/os2.pdf
//
// ARC policy:
//	B1 <- [ T1 <- -> T2 ] -> B2
//	L1 = B1 + T1, L2 = B2 + T2, 0 ≤ |L1| + |L2| ≤ 2c
//	0 ≤ |L1| ≤ c , L2 can be bigger than c
//	When a page is accessed
//		if the page is in L1 ∪ L2 (union), move it to the MRU of L2
//		otherwise move it to the MRU of L1
//	If adding the new page makes |L1| + |L2| > 2c or |L1| > c
//		if L1 (before the addition) contains less than c pages, take out the LRU of L2
//		otherwise take out the LRU of L1
type ARC struct {
	size int

	pairs map[interface{}]*item

	b1, t1, t2, b2 *arcList

	p int // for partition

	onEvict EvictCallback
}

type item struct {
	key   interface{}
	value interface{}
}

// EvictCallback called when a lru entry is evicted
type EvictCallback func(key interface{}, value interface{})

// NewARC creates a new ARC cache with input size and onEvict function
func NewARC(size int, onEvict EvictCallback) (*ARC, error) {
	if size <= 0 {
		return nil, errors.New("invalid cache size")
	}
	return &ARC{
		size:    size,
		pairs:   make(map[interface{}]*item),
		b1:      newArcList(),
		t1:      newArcList(),
		t2:      newArcList(),
		b2:      newArcList(),
		onEvict: onEvict,
	}, nil
}

// Add adds key value to the cache and returns true if an eviction occurred
func (arc *ARC) Add(key, value interface{}) (found, evicted bool) {
	// if found, update value
	if it, found := arc.pairs[key]; found {
		it.value = value
		arc.updateKey(key)
		return true, false
	}

	// not found, add new item
	it := &item{
		key:   key,
		value: value,
	}
	arc.pairs[key] = it
	evicted = arc.addNewItem(it)
	return
}

// Get returns value if key and update timestamp
func (arc *ARC) Get(key interface{}) (value interface{}, found bool) {
	it, found := arc.pairs[key]
	if found {
		value = it.value
		arc.move(key)
	}
	return
}

// Peek returns value of key without updating timestamp
func (arc *ARC) Peek(key interface{}) (value interface{}, found bool) {
	it, found := arc.pairs[key]
	if found {
		value = it.value
	}
	return
}

// Contains returns true if key found in the cache without updating timestamp
func (arc *ARC) Contains(key interface{}) (found bool) {
	_, found = arc.pairs[key]
	return
}

// Remove removes key value if found in the cache
func (arc *ARC) Remove(key interface{}) (value interface{}, found bool) {
	// not in memory (t1 + t2)
	it, found := arc.pairs[key]
	if !found {
		return
	}
	value = it.value

	// in memory (t1 + t2)
	if e := arc.t1.GetElement(key); e != nil {
		found = true
		arc.t1.Remove(key, e)
		delete(arc.pairs, key)
		arc.b1.PushFront(key)
		if arc.onEvict != nil {
			arc.onEvict(key, value)
		}
		return
	}
	if e := arc.t2.GetElement(key); e != nil {
		found = true
		arc.t2.Remove(key, e)
		delete(arc.pairs, key)
		arc.b2.PushFront(key)
		if arc.onEvict != nil {
			arc.onEvict(key, value)
		}
		return
	}
	return
}

// Size returns key value pair numbers in the cache
func (arc *ARC) Size() int {
	return len(arc.pairs)
}

// Resize resize cache size and returns diff
func (arc *ARC) Resize(size int) (diff int, err error) {
	diff = size - arc.size
	if size > 0 {
		arc.size = size
	} else {
		return 0, errors.New("resize size should be larger than 0")
	}

	// extend
	if diff >= 0 {
		return
	}

	// diff < 0 => shrink
	for i := 0; i < -diff; i++ {
		// arc.removeLeastUsed()
		arc.RemoveLeastUsed()
	}
	return
}

// GetLeastUsed returns least used key value pairs if found in the cache
// not update timestamp
// t1 -> t2
func (arc *ARC) GetLeastUsed() (key, value interface{}, found bool) {
	var e *list.Element
	if arc.t1.Len() > 0 {
		e = arc.t1.l.Back()
		key = e.Value
		value = arc.pairs[key].value
		found = true
		return
	}
	if arc.t2.Len() > 0 {
		e = arc.t2.l.Back()
		key = e.Value
		value = arc.pairs[key].value
		found = true
		return
	}
	return
}

// RemoveLeastUsed removes and returns least used key value pairs if found in the cache
// t1 -> t2
func (arc *ARC) RemoveLeastUsed() (key, value interface{}, found bool) {
	var e *list.Element
	if arc.t1.Len() > 0 {
		e = arc.t1.l.Back()
		key = e.Value
		value, found = arc.Remove(key)
		return
	}
	if arc.t2.Len() > 0 {
		e = arc.t2.l.Back()
		key = e.Value
		value, found = arc.Remove(key)
		return
	}
	return
}

// Clear clears all pairs in the cache
func (arc *ARC) Clear() {
	for _, v := range arc.pairs {
		if arc.onEvict != nil {
			arc.onEvict(v.key, v.value)
		}
	}
	arc.pairs = make(map[interface{}]*item)
	arc.t1 = newArcList()
	arc.b1 = newArcList()
	arc.t1 = newArcList()
	arc.b2 = newArcList()
}

func (arc *ARC) printPairs() map[interface{}]interface{} {
	res := make(map[interface{}]interface{}, arc.Size())
	for k, v := range arc.pairs {
		res[k] = v.value
	}
	return res
}

// move key from t1 to t2 or from t2 to t2's front
func (arc *ARC) move(key interface{}) {
	if e := arc.t1.GetElement(key); e != nil {
		// move key from t1 to t2
		arc.t1.Remove(key, e)
		arc.t2.PushFront(key)
		return
	}

	if e := arc.t2.GetElement(key); e != nil {
		// move key from t2 to t2 front (MRU)
		arc.t2.MoveFront(e)
		return
	}
}

// if the page is in L1 ∪ L2 (valid cache c)
// update key and set p (when arc is full)
func (arc *ARC) updateKey(key interface{}) {
	// The increments and the decrements are subject to the stipulation 0 ≤ p ≤ c

	// If there is a hit in T1 or T2, do nothing
	if arc.t1.Contains(key) || arc.t2.Contains(key) {
		return
	}
	// If there is a hit in B1
	//	If the size of B1 is at least the size of B2, increment p by 1
	//	otherwise, increment p by |B2|-|B1|
	if e := arc.b1.GetElement(key); e != nil {
		// set p
		if arc.isFull() {
			if arc.b1.Len() >= arc.b2.Len() {
				arc.p++
			} else {
				arc.p += arc.b2.Len() - arc.b1.Len()
			}
		}
		// move key from b1 to t2
		arc.moveOutMemory(key)
		arc.b1.Remove(key, e)
		arc.t2.PushFront(key)
		return
	}

	// If there is a hit in B2
	//	If the size of B2 is at least the size of B1, decrement p by 1
	//	otherwise, decrement p by |B1|-|B2|
	if e := arc.b2.GetElement(key); e != nil {
		// set p
		if arc.isFull() {
			if arc.b2.Len() >= arc.b1.Len() {
				arc.p--
			} else {
				arc.p -= arc.b1.Len() - arc.b2.Len()
			}
		}
		// move key from b2 to t2
		arc.moveOutMemory(key)
		arc.b2.Remove(key, e)
		arc.t2.PushFront(key)
		return
	}
}

// if adding the new page makes |L1| + |L2| > 2c or |L1| > c
//	if L1 (before the addition) contains less than c pages, take out the LRU of L2
//	otherwise take out the LRU of L1
func (arc *ARC) addNewItem(it *item) (evicted bool) {
	// |L1| + |L2| > 2c or |L1| > c
	if arc.isFull() && arc.t1.Len()+arc.b1.Len() == arc.size {
		evicted = true
		// if L1 (before the addition) contains less than c pages, take out the LRU of L2
		if arc.t1.Len() < arc.size {
			arc.b1.RemoveBack()
			arc.moveOutMemory(it.key)
			// otherwise take out the LRU of L1
		} else {
			lruKey := arc.t1.RemoveBack()
			if item, found := arc.pairs[lruKey]; found {
				delete(arc.pairs, lruKey)
				if arc.onEvict != nil {
					arc.onEvict(item.key, item.value)
				}
			}
		}
	} else {
		total := arc.t1.Len() + arc.b1.Len() + arc.t2.Len() + arc.b2.Len()
		// check evict
		if total >= arc.size {
			if total == (2 * arc.size) {
				if arc.b2.Len() > 0 {
					arc.b2.RemoveBack()
				} else {
					arc.b1.RemoveBack()
				}
			}
			arc.moveOutMemory(it.key)
			evicted = true
		}
	}
	// add new item to t1
	arc.t1.PushFront(it.key)
	return
}

// From http://u.cs.biu.ac.il/~wiseman/2os/2os/os2.pdf
// When a page is moved from a "T" list to a "B" list, it will be taken out of the memory
// Let p be the current target size for the list T1
// If |T1| > p, move the LRU of T1 to be the MRU of B1
// If |T1| < p, move the LRU of T2 to be the MRU of B2
// If |T1| = p,
//	– If the accessed page has been in B2, move the LRU of T1 to be the MRU of B1 (Because p is going to be decremented)
//	– If the accessed page has been in B1 or has not been in the memory, move the LRU of T2 to be the MRU of B2
func (arc *ARC) moveOutMemory(key interface{}) {
	if !arc.isFull() {
		return
	}

	var movedKey interface{}
	// |T1| > p || ((|T1| = p) && the accessed page has been in B2)
	if arc.t1.Len() > arc.p || (arc.t1.Len() > 0 && (arc.t1.Len() == arc.p && arc.b2.Contains(key))) {
		// move the LRU of T1
		movedKey = arc.t1.RemoveBack()
		// to be the MRU of B1
		arc.b1.PushFront(movedKey)
		// |T1| < p || ((|T1| = p) && the accessed page has been in B1 or has not been in the memory)
	} else if arc.t2.Len() > 0 {
		// move the LRU of T2
		movedKey = arc.t2.RemoveBack()
		// to be the MRU of B2
		arc.b2.PushFront(movedKey)
	} else {
		// move the LRU of T1
		movedKey = arc.t1.RemoveBack()
		// to be the MRU of B1
		arc.b1.PushFront(movedKey)
	}

	// evict movedKey's item
	if item, found := arc.pairs[movedKey]; found {
		delete(arc.pairs, movedKey)
		if arc.onEvict != nil {
			arc.onEvict(item.key, item.value)
		}
	}
}

func (arc *ARC) isFull() bool {
	return (arc.t1.Len() + arc.t2.Len()) == arc.size
}

type arcList struct {
	l    *list.List
	keys map[interface{}]*list.Element
}

func newArcList() *arcList {
	return &arcList{
		l:    list.New(),
		keys: make(map[interface{}]*list.Element),
	}
}

func (al *arcList) Contains(key interface{}) (found bool) {
	_, found = al.keys[key]
	return
}

func (al *arcList) GetElement(key interface{}) *list.Element {
	return al.keys[key]
}

func (al *arcList) PushFront(key interface{}) {
	if e, found := al.keys[key]; found {
		al.l.MoveToFront(e)
		return
	}
	e := al.l.PushFront(key)
	al.keys[key] = e
}

func (al *arcList) MoveFront(e *list.Element) {
	al.l.MoveToFront(e)
}

func (al *arcList) Remove(key interface{}, e *list.Element) {
	delete(al.keys, key)
	al.l.Remove(e)
}

func (al *arcList) RemoveBack() (key interface{}) {
	e := al.l.Back()
	key = e.Value

	al.l.Remove(e)
	delete(al.keys, key)
	return
}

func (al *arcList) Len() int {
	return al.l.Len()
}
