package bidirectionalmap

import "goContainer/maps"

// BidirectionalMap struct
//	simply based on two maps
type BidirectionalMap struct {
	flip, flop map[interface{}]interface{}
}

// NewBidirectionalMap creates a new bidirectional map
func NewBidirectionalMap() *BidirectionalMap {
	return &BidirectionalMap{
		flip: make(map[interface{}]interface{}),
		flop: make(map[interface{}]interface{}),
	}
}

// Set sets key value pairs
func (m *BidirectionalMap) Set(key, value interface{}) {
	m.flip[key] = value
	m.flop[value] = key
}

// Get gets value with input key if found inside map
func (m *BidirectionalMap) Get(key interface{}) (value interface{}, found bool) {
	if value, found = m.flip[key]; found {
		return
	} else if value, found = m.flop[key]; found {
		return
	}
	return
}

// Delete deletes key value pairs
func (m *BidirectionalMap) Delete(key interface{}) bool {
	v1, found1 := m.flip[key]
	v2, found2 := m.flop[key]
	if found1 {
		delete(m.flip, key)
		delete(m.flop, v1)
		return true
	}
	if found2 {
		delete(m.flip, v2)
		delete(m.flop, key)
		return true
	}
	return false
}

// Empty returns true if no kv pairs inside map
func (m *BidirectionalMap) Empty() bool {
	return len(m.flip) == 0
}

// Size returns quantity of kv pairs inside map
func (m *BidirectionalMap) Size() int {
	return len(m.flip)
}

// Clear clears the map
func (m *BidirectionalMap) Clear() {
	*m = *NewBidirectionalMap()
}

// Values returns all values un-ordered
func (m *BidirectionalMap) Values() []interface{} {
	values := make([]interface{}, 0, len(m.flip))
	for _, v := range m.flip {
		values = append(values, v)
	}
	return values
}

// Keys returns all keys un-ordered
func (m *BidirectionalMap) Keys() []interface{} {
	keys := make([]interface{}, 0, len(m.flip))
	for k := range m.flip {
		keys = append(keys, k)
	}
	return keys
}

type iterator struct {
	m    *map[interface{}]interface{}
	keys []interface{}
}

// HasNext to meet Iterator interface
func (it *iterator) HasNext() bool {
	return len(it.keys) != 0
}

// Next to meet Iterator interface
func (it *iterator) Next() (key, value interface{}) {
	key = it.keys[0]
	value = (*it.m)[key]
	it.keys = it.keys[1:]
	return
}

// Iterator returns iterator
func (m *BidirectionalMap) Iterator() maps.Iterator {
	return &iterator{
		m:    &m.flip,
		keys: m.Keys(),
	}
}
