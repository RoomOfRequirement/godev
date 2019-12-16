package eventloop

import (
	"context"
	"time"
)

// EventLoop interface ...
type EventLoop interface {
	Push(event *Event)
	Emit(eventName string, data interface{})
	On(eventName string, callback Callback) error
	OnWithTimeout(eventName string, callback Callback, timeout time.Duration) error
	Start()
	Stop()
}

// Event struct ...
type Event struct {
	Name string
	Data interface{}
}

// Callback ...
type Callback func(ctx context.Context, args ...interface{}) error
