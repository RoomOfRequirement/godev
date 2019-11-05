package linkedhashmap

import (
	"container/list"
	"goContainer/maps"
)

type item struct {
	key   *list.Element
	value interface{}
}

// LinkedHashMap struct uses a linked list to record key insertion order
type LinkedHashMap struct {
	_list *list.List
	_map  map[interface{}]*item
}

// NewLinkedHashMap returns a new linked hash map
func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		_list: list.New(),
		_map:  make(map[interface{}]*item),
	}
}

// Set sets key value pairs
func (m *LinkedHashMap) Set(key, value interface{}) {
	if _, found := m._map[key]; !found {
		e := m._list.PushBack(key)
		m._map[key] = &item{
			key:   e,
			value: value,
		}
		return
	}
	m._map[key].value = value
}

// Get gets value with input key if found inside map
func (m *LinkedHashMap) Get(key interface{}) (value interface{}, found bool) {
	it, found := m._map[key]
	value = it.value
	return
}

// Delete deletes key value pairs
func (m *LinkedHashMap) Delete(key interface{}) bool {
	if it, found := m._map[key]; found {
		delete(m._map, key)
		m._list.Remove(it.key)
		return true
	}
	return false
}

// Empty returns true if no kv pairs inside map
func (m *LinkedHashMap) Empty() bool {
	return len(m._map) == 0
}

// Size returns quantity of kv pairs inside map
func (m *LinkedHashMap) Size() int {
	return len(m._map)
}

// Clear clears the map
func (m *LinkedHashMap) Clear() {
	*m = *NewLinkedHashMap()
}

// Values returns all values in insertion order
func (m *LinkedHashMap) Values() []interface{} {
	values := make([]interface{}, 0, len(m._map))
	for e := m._list.Front(); e != nil; e = e.Next() {
		values = append(values, m._map[e.Value].value)
	}
	return values
}

// Keys returns all keys in insertion order
func (m *LinkedHashMap) Keys() []interface{} {
	keys := make([]interface{}, 0, len(m._map))
	for e := m._list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value)
	}
	return keys
}

type iterator struct {
	m      *LinkedHashMap
	cursor *list.Element
}

// HasNext to meet Iterator interface
func (iter *iterator) HasNext() bool {
	return iter.cursor.Next() != nil
}

// Next to meet Iterator interface
func (iter *iterator) Next() (key, value interface{}) {
	it := iter.m._map[iter.cursor.Value]
	iter.cursor = iter.cursor.Next()
	return it.key.Value, it.value
}

// Iterator returns iterator, details can be found in red-black tree implementation
func (m *LinkedHashMap) Iterator() maps.Iterator {
	iter := iterator{
		m:      m,
		cursor: m._list.Front(),
	}
	return &iter
}
