package timewheel

import (
	"fmt"
	"godev/controlflow/logger"
	"godev/utils"
)

// Logger ... global logger
var Logger = logger.NewLogger("info")

// Task wrapper
type Task struct {
	fn   func(...interface{}) // task func
	args []interface{}        // args in order
}

// NewTask ...
func NewTask(fn func(...interface{}), args []interface{}) *Task {
	return &Task{
		fn:   fn,
		args: args,
	}
}

// Run ...
func (t *Task) Run() {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error(utils.Trace(fmt.Sprintf("%s", err)))
		}
	}()
	t.fn(t.args...)
}
