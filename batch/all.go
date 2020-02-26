package batch

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"strings"
	"sync"
)

// All ...
type All struct {
	cancel  func()
	wg      *sync.WaitGroup
	sema    *semaphore.Weighted
	semaCtx context.Context
	errOnce *sync.Once
	errs    chan error
}

// NewAll ...
//	executes n task without concurrency limit
//	number of errors: [0, n]
func NewAll(taskN int) *All {
	if taskN < 1 {
		taskN = 1
	}
	return &All{
		cancel:  nil,
		wg:      &sync.WaitGroup{},
		sema:    nil,
		semaCtx: nil,
		errOnce: nil,
		errs:    make(chan error, taskN),
	}
}

// NewAllWithLimit ...
//	executes n task with concurrency limit: max concurrency goroutines
//	number of errors: [0, n]
func NewAllWithLimit(taskN, concurrency int) *All {
	if taskN < 1 {
		taskN = 1
	}
	if concurrency < 1 {
		concurrency = 1
	}
	return &All{
		cancel:  nil,
		wg:      &sync.WaitGroup{},
		sema:    semaphore.NewWeighted(int64(concurrency)), // equally weighted
		semaCtx: context.TODO(),
		errOnce: nil,
		errs:    make(chan error, taskN),
	}
}

// WithContext ...
//	this will replace semaCtx with ctx,
//	which means a ctx canceled error will be recorded when one of task throws error
func (a *All) WithContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	a.cancel = cancel
	a.errOnce = &sync.Once{}
	// replace default semaCtx
	a.semaCtx = ctx
	return ctx
}

// Go ...
func (a *All) Go(idx int, f func() error) {
	a.wg.Add(1)
	if a.sema != nil {
		// err when semaCtx is Done (semaCtx.Err())
		err := a.sema.Acquire(a.semaCtx, 1)
		if err != nil {
			a.errs <- fmt.Errorf("func %d return with error: %s", idx, err)
			a.wg.Done()
			return
		}
	}
	go func() {
		defer a.wg.Done()
		if a.sema != nil {
			defer a.sema.Release(1)
		}
		if err := f(); err != nil {
			if a.cancel != nil {
				a.errOnce.Do(a.cancel)
			}
			a.errs <- fmt.Errorf("func %d return with error: %s", idx, err)
		}
	}()
}

// Wait ...
func (a *All) Wait() error {
	a.wg.Wait()
	close(a.errs)
	if len(a.errs) == 0 {
		return nil
	}
	var errStr strings.Builder
	for err := range a.errs {
		errStr.WriteString(fmt.Sprintf("%s\n", err))
	}
	return fmt.Errorf("%s", errStr.String())
}
