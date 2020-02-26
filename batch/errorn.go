package batch

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// ErrorN allows n error
//	like `x/sync/errgroup` but allow n errors
type ErrorN struct {
	cancel  func()
	wg      *sync.WaitGroup
	errOnce *sync.Once
	errs    chan error
}

// NewErrorN ... errorN at least 1
func NewErrorN(errorN int) *ErrorN {
	if errorN < 1 {
		errorN = 1
	}
	return &ErrorN{
		cancel:  nil,
		wg:      &sync.WaitGroup{},
		errOnce: nil,
		errs:    make(chan error, errorN),
	}
}

// WithContext ... cancel context once error occurs
func (en *ErrorN) WithContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	en.cancel = cancel
	en.errOnce = &sync.Once{}
	return ctx
}

// Go ...
//	notice: this may execute more than N tasks when runtime.NumCPU > N
//	if you want to control the concurrency, you may use semaphore
func (en *ErrorN) Go(idx int, f func() error) {
	en.wg.Add(1)
	go func() {
		defer en.wg.Done()
		// full
		if len(en.errs) == cap(en.errs) {
			return
		}
		if err := f(); err != nil {
			select {
			case en.errs <- fmt.Errorf("func %d return with error: %s", idx, err):
				// push in err
				if en.cancel != nil {
					en.errOnce.Do(en.cancel)
				}
			default:
				// full
				return
			}
		}
	}()
}

// Wait ...
func (en *ErrorN) Wait() error {
	en.wg.Wait()
	close(en.errs)
	if len(en.errs) == 0 {
		return nil
	}
	var errStr strings.Builder
	for err := range en.errs {
		errStr.WriteString(fmt.Sprintf("%s\n", err))
	}
	return fmt.Errorf("%s", errStr.String())
}
