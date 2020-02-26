package batch

import "context"

// Batch interface ...
type Batch interface {
	// WithContext ... context with cancel, it will be cancelled once error occurs
	WithContext(ctx context.Context) context.Context
	// Go ... call in for loop, pass in loop idx
	Go(idx int, f func() error)
	// Wait ...
	Wait() error
}
