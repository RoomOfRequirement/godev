package flow

import (
	"context"
	"sync"
	"time"
)

type throttler struct {
	lock  sync.Mutex
	wait  time.Duration
	timer *time.Timer
}

func (t *throttler) throttledLeading(ctx context.Context, f func()) {
	t.lock.Lock()
	defer t.lock.Unlock()
	select {
	case <-ctx.Done():
		// cancel
	default:
		// call on leading
		if t.timer == nil {
			f()
			t.timer = time.AfterFunc(t.wait, func() {
				t.lock.Lock()
				t.timer.Stop()
				t.timer = nil
				t.lock.Unlock()
			})
		}
		// wait for timer == nil
	}
}

func (t *throttler) throttledTailing(ctx context.Context, f func()) {
	t.lock.Lock()
	defer t.lock.Unlock()
	select {
	case <-ctx.Done():
		// cancel
	default:
		// call on tailing
		if t.timer == nil {
			t.timer = time.AfterFunc(t.wait, func() {
				f()
				t.lock.Lock()
				t.timer.Stop()
				t.timer = nil
				t.lock.Unlock()
			})
		}
		// wait for timer == nil
	}
}

// NewThrottle ...
//	f can only be executed once in `wait` time interval, no matter leading or not
//	reference: https://lodash.com/docs/#throttle
func NewThrottle(ctx context.Context, wait time.Duration, leading bool) (throttle func(f func()), cancel func()) {
	t := &throttler{
		lock:  sync.Mutex{},
		wait:  wait,
		timer: nil,
	}
	ctx, cancel = context.WithCancel(ctx)
	if leading {
		return func(f func()) {
			t.throttledLeading(ctx, f)
		}, cancel
	}
	return func(f func()) {
		t.throttledTailing(ctx, f)
	}, cancel
}
