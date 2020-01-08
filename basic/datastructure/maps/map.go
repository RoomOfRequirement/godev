package maps

import (
	"goContainer/basic"
)

// Map interface
type Map interface {
	Set(key, value interface{})
	Get(key interface{}) (value interface{}, found bool)
	Delete(key interface{}) bool
	Keys() []interface{}
	Iterator() Iterator

	basic.Container
}

// Iterator interface
type Iterator interface {
	HasNext() bool
	Next() (key interface{}, value interface{})
}
