package batch

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// TaskN completes at least n task and records errors
//	number of errors can be [0, n]
type TaskN struct {
	cancel  func()
	wg      *sync.WaitGroup
	errOnce *sync.Once
	errs    chan error
	n       int32
}

// NewTaskN ... taskN at least 1
func NewTaskN(taskN int) *TaskN {
	if taskN < 1 {
		taskN = 1
	}
	return &TaskN{
		cancel:  nil,
		wg:      &sync.WaitGroup{},
		errOnce: nil,
		errs:    make(chan error, taskN),
		n:       int32(taskN),
	}
}

// WithContext ...
func (tn *TaskN) WithContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	tn.cancel = cancel
	tn.errOnce = &sync.Once{}
	return ctx
}

// Go ...
//	notice: this may execute more than N tasks when runtime.NumCPU > N
//	if you want to control the concurrency, you may use semaphore
func (tn *TaskN) Go(idx int, f func() error) {
	tn.wg.Add(1)
	go func() {
		defer tn.wg.Done()
		if atomic.LoadInt32(&tn.n) == 0 {
			return
		}
		// full
		if len(tn.errs) == cap(tn.errs) {
			return
		}
		if err := f(); err != nil {
			select {
			case tn.errs <- fmt.Errorf("func %d return with error: %s", idx, err):
				// push in err
				if tn.cancel != nil {
					tn.errOnce.Do(tn.cancel)
				}
			default:
				// full
				return
			}
		} else {
			// full
			if atomic.LoadInt32(&tn.n) == 0 {
				return
			}
			atomic.AddInt32(&tn.n, -1)
		}
	}()
}

// Wait ...
func (tn *TaskN) Wait() error {
	tn.wg.Wait()
	close(tn.errs)
	if len(tn.errs) == 0 {
		return nil
	}
	var errStr strings.Builder
	for err := range tn.errs {
		errStr.WriteString(fmt.Sprintf("%s\n", err))
	}
	return fmt.Errorf("%s", errStr.String())
}
