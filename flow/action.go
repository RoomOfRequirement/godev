package flow

import (
	"context"
)

// Action interface
type Action interface {
	Execute(ctx context.Context) error
}

// ActionFunc meets Action interface and makes a function into action
type ActionFunc func(ctx context.Context) error

// Execute function to meet Action / Executor interface
func (af ActionFunc) Execute(ctx context.Context) error {
	return af(ctx)
}
