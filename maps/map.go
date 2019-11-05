package maps

import container "goContainer"

// Map interface
type Map interface {
	Set(key, value interface{})
	Get(key interface{}) (value interface{}, found bool)
	Delete(key interface{}) bool
	Keys() []interface{}
	Iterator() Iterator

	container.Container
}

// Iterator interface
type Iterator interface {
	HasNext() bool
	Next() (key interface{}, value interface{})
}
