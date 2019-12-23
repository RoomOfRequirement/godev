package flow

// reference: https://davidwalsh.name/javascript-debounce-function

import (
	"sync"
	"time"
)

// NewDebounce wraps a function and returns a wrapped function, that,
// as long as it continues to be invoked, will not be triggered
// The function will be called after it stops being called for N milliseconds
// If `immediate` is passed, trigger the function on the leading edge, instead of the trailing
func NewDebounce(after time.Duration, immediate bool) (debounce func(f func())) {
	d := &debouncer{
		after: after,
	}
	if !immediate {
		return func(f func()) {
			d.debounce(f)
		}
	}
	return func(f func()) {
		d.debounced(f)
	}
}

type debouncer struct {
	lock  sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (db *debouncer) debounce(f func()) {
	db.lock.Lock()

	// if called, replace current timer with new one
	if db.timer != nil {
		db.timer.Stop()
	}
	db.timer = time.AfterFunc(db.after, f)

	db.lock.Unlock()
}

func (db *debouncer) debounced(f func()) {
	db.lock.Lock()

	if db.timer == nil {
		f()
		db.timer = time.AfterFunc(db.after, func() {
			db.lock.Lock()
			db.timer.Stop()
			db.timer = nil
			db.lock.Unlock()
		})
	} else {
		db.timer.Stop()
		db.timer = time.AfterFunc(db.after, func() {
			db.lock.Lock()
			db.timer.Stop()
			db.timer = nil
			db.lock.Unlock()
		})
	}

	db.lock.Unlock()
}
