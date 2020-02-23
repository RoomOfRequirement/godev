package timewheel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTask(t *testing.T) {
	a, b, c := 0, 1, 2
	task := NewTask(func(i ...interface{}) {
		*i[0].(*int), *i[1].(*int), *i[2].(*int) = 2, 1, 0
	}, []interface{}{&a, &b, &c})
	task.Run()
	assert.Equal(t, 2, a)
	assert.Equal(t, 1, b)
	assert.Equal(t, 0, c)

	// panic
	task = NewTask(func(i ...interface{}) {
		panic("hi, panic here")
	}, []interface{}{})
	assert.NotPanics(t, task.Run)
}
