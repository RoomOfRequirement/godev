package trylock

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// see: https://golang.org/src/sync/mutex.go
const mutexLocked = 1 << iota // mutex is locked

type tryMutex struct {
	sync.Mutex
}

func (m *tryMutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}

func (m *tryMutex) TryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			timer.Stop()
			return false
		default:
			if m.TryLock() {
				return true
			}
		}
	}
}

type tryRWMutex struct {
	// TODO: better to divide state into two?
	state   int32         // 0: no lock, -1: write lock, >=1: read lock
	sig     chan struct{} // like sync.Cond but simpler
	sigLock sync.Mutex
}

func (rw *tryRWMutex) TryLock() bool {
	if atomic.CompareAndSwapInt32(&rw.state, 0, -1) {
		return true
	}
	return false
}

func (rw *tryRWMutex) TryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	rw.sigLock.Lock()
	sig := rw.sig
	rw.sigLock.Unlock()
	for {
		select {
		case <-timer.C:
			timer.Stop()
			return false
		case <-sig:
			if rw.TryLock() {
				return true
			}
		}
	}
}

func (rw *tryRWMutex) RTryLock() bool {
	if state := atomic.LoadInt32(&rw.state); state >= 0 {
		if atomic.CompareAndSwapInt32(&rw.state, state, state+1) {
			return true
		}
	}
	return false
}

func (rw *tryRWMutex) RTryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	rw.sigLock.Lock()
	sig := rw.sig
	rw.sigLock.Unlock()
	for {
		select {
		case <-timer.C:
			timer.Stop()
			return false
		case <-sig:
			if rw.RTryLock() {
				return true
			}
		}
	}
}

func (rw *tryRWMutex) Lock() {
	rw.sigLock.Lock()
	sig := rw.sig
	rw.sigLock.Unlock()
	for {
		if rw.TryLock() {
			return
		}
		select {
		case <-sig:
			if rw.TryLock() {
				return
			}
		}
	}
}

func (rw *tryRWMutex) Unlock() {
	// TODO: can this failed?
	atomic.CompareAndSwapInt32(&rw.state, -1, 0)

	newSig := make(chan struct{}, 1)
	rw.sigLock.Lock()
	sig := rw.sig
	rw.sig = newSig
	rw.sigLock.Unlock()

	// broadcast
	close(sig)
}

func (rw *tryRWMutex) RLock() {
	rw.sigLock.Lock()
	sig := rw.sig
	rw.sigLock.Unlock()
	for {
		if rw.RTryLock() {
			return
		}
		select {
		case <-sig:
			if rw.RTryLock() {
				return
			}
		}
	}
}

func (rw *tryRWMutex) RUnlock() {
	state := atomic.AddInt32(&rw.state, -1)
	if state < 0 {
		panic("RUnlock() failed")
	}

	// free for others
	if state == 0 {
		newSig := make(chan struct{}, 1)
		rw.sigLock.Lock()
		sig := rw.sig
		rw.sig = newSig
		rw.sigLock.Unlock()

		// broadcast
		close(sig)
	}
}
