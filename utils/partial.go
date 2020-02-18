package utils

import (
	"errors"
	"reflect"
)

// PartialFunc based on reflect
func PartialFunc(funcMap map[string]interface{}, funcName string, funcArgs ...interface{}) (resSlice []reflect.Value, err error) {
	f := reflect.ValueOf(funcMap[funcName])
	if len(funcArgs) != f.Type().NumIn() {
		err = errors.New("invalid number of funcArgs")
		return
	}
	in := make([]reflect.Value, len(funcArgs))
	for idx, arg := range funcArgs {
		in[idx] = reflect.ValueOf(arg)
	}
	resSlice = f.Call(in)
	return
}
