package spinlock

import (
	"runtime"
	"sync/atomic"
	"time"
)

// Lock ...
type Lock struct {
	p uintptr
}

// TryLock ...
func (l *Lock) TryLock() bool {
	return atomic.CompareAndSwapUintptr(&l.p, 0, 1)
}

// IsLocked ...
func (l *Lock) IsLocked() bool {
	return atomic.LoadUintptr(&l.p) == 1
}

// TryUnlock ...
func (l *Lock) TryUnlock() bool {
	return atomic.CompareAndSwapUintptr(&l.p, 1, 0)
}

// ForceUnlock ...
func (l *Lock) ForceUnlock() {
	atomic.StoreUintptr(&l.p, 0)
}

// SpinLock ...
func (l *Lock) SpinLock(retryInterval, timeout time.Duration) bool {
	end := time.Now().Add(timeout)
	for {
		if l.TryLock() {
			return true
		} else if time.Now().After(end) {
			return false
		}
		time.Sleep(retryInterval)
	}
}

// ForceLock ...
func (l *Lock) ForceLock() {
	if !l.TryLock() {
		// yields the processor but not suspend current goroutine,
		// current execution will resume automatically
		runtime.Gosched()
	}
}
