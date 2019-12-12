package flow

import (
	"context"
	"errors"
	"goContainer/pool/goroutinepool"
	"golang.org/x/sync/errgroup"
)

// Executor interface
type Executor interface {
	Execute(ctx context.Context, actions ...Action) error
}

// SequentialExecutor implements a sequential executor
type SequentialExecutor struct{}

// Execute function to meet Executor interface
func (SequentialExecutor) Execute(ctx context.Context, actions ...Action) error {
	for _, a := range actions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := a.Execute(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// ConcurrentExecutor implements a concurrent executor, every action is executed in a new goroutine
type ConcurrentExecutor struct{}

// Execute function to meet Executor interface
func (ConcurrentExecutor) Execute(ctx context.Context, actions ...Action) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// doc: https://godoc.org/golang.org/x/sync/errgroup
	grp, ctx := errgroup.WithContext(ctx)

	for _, a := range actions {
		grp.Go(goFunc(ctx, a))
	}

	return grp.Wait()
}

// wrap action.Execute to meet grp.Go func argument signature
func goFunc(ctx context.Context, a Action) func() error {
	return func() error {
		return a.Execute(ctx)
	}
}

// PoolExecutor implements a pool executor, actions are executed in a goroutine pool
type PoolExecutor struct {
	p *goroutinepool.Pool
}

// ErrPoolClosed ...
var ErrPoolClosed = errors.New("pool closed")

// NewPool creates a new goroutine pool executor
func NewPool(workerNum int) (Pool Executor, StopFunc func()) {
	p := goroutinepool.New(workerNum)
	Pool = PoolExecutor{p: p}
	StopFunc = p.Stop
	return
}

// Execute function to meet Executor interface
func (pe PoolExecutor) Execute(ctx context.Context, actions ...Action) error {
	actionNum := len(actions)
	if actionNum == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res := make(chan error, actionNum)

	var err error
	var queued int

Loop:
	for _, a := range actions {
		if pe.p.Stopped() {
			cancel()
			return ErrPoolClosed
		}
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break Loop
		default:
			pe.p.Submit(func() {
				res <- a.Execute(ctx)
			})
			queued++
		}
	}

	for ; queued > 0; queued-- {
		if r := <-res; r != nil {
			if err == nil {
				err = r
				cancel()
			}
		}
	}

	return err
}
