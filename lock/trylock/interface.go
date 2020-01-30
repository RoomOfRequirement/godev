package trylock

import (
	"sync"
	"time"
)

// TryMutex wraps Mutex with additional `try` methods
type TryMutex interface {
	TryLock() bool
	TryLockWithTimeout(timeout time.Duration) bool
	Lock()
	Unlock()
}

// TryRWMutex realizes RWMutex with `try` methods
type TryRWMutex interface {
	TryMutex
	RTryLock() bool
	RTryLockWithTimeout(timeout time.Duration) bool
	RLock()
	RUnlock()
}

// NewTryMutex returns a TryMutex
func NewTryMutex() TryMutex {
	return &tryMutex{
		Mutex: sync.Mutex{},
	}
}

// NewTryRWMutex returns a TryRWMutex
func NewTryRWMutex() TryRWMutex {
	return &tryRWMutex{
		state:   0,
		sig:     make(chan struct{}, 1),
		sigLock: sync.Mutex{},
	}
}
