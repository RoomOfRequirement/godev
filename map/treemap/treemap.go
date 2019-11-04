package treemap

import (
	container "goContainer"
	"goContainer/tree/rbtree"
)

// Map based on red-black tree
type Map struct {
	tree *rbtree.RBTree
}

// NewMap creates a new map with input comparator
func NewMap(comparator container.Comparator) *Map {
	m := &Map{
		tree: new(rbtree.RBTree),
	}
	m.tree.Comparator = comparator
	return m
}

// Set sets key value pairs
func (m *Map) Set(key, value interface{}) {
	m.tree.Update(key, value)
}

// Get gets value with input key if found inside map
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	return m.tree.Get(key)
}

// Delete deletes key value pairs
func (m *Map) Delete(key interface{}) bool {
	return m.tree.Delete(key)
}

// Keys returns all keys
func (m *Map) Keys() []interface{} {
	return m.tree.Keys()
}

// Values returns all values
func (m *Map) Values() []interface{} {
	return m.tree.Values()
}

// Empty returns true if no kv pairs inside map
func (m *Map) Empty() bool {
	return m.tree.Empty()
}

// Size returns quantity of kv pairs inside map
func (m *Map) Size() int {
	return m.tree.Size()
}

// Clear clears the map
func (m *Map) Clear() {
	m.tree.Clear()
}

// Iterator returns iterator, details can be found in red-black tree implementation
func (m *Map) Iterator() *rbtree.Iterator {
	return m.tree.Iterator()
}
