package trylock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTryMutex(t *testing.T) {
	m := NewTryMutex()
	m.Lock()
	assert.False(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
	m.Unlock()

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.True(t, m.TryLockWithTimeout(15*time.Millisecond))
	m.Unlock()

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.False(t, m.TryLockWithTimeout(5*time.Millisecond))
}

func TestNewTryRWMutex(t *testing.T) {
	m := NewTryRWMutex()
	m.Lock()
	assert.False(t, m.TryLock())
	assert.False(t, m.RTryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
	assert.False(t, m.RTryLock())
	m.Unlock()
	m.RLock()
	assert.True(t, m.RTryLock())
	assert.True(t, m.RTryLock())
	m.RUnlock()
	m.RUnlock()
	m.RUnlock()

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.True(t, m.TryLockWithTimeout(15*time.Millisecond))
	m.Unlock()

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.False(t, m.TryLockWithTimeout(5*time.Millisecond))

	m.RLock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.RUnlock()
	})
	assert.True(t, m.RTryLockWithTimeout(15*time.Millisecond))
	m.RUnlock()

	m.RLock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.RUnlock()
	})
	assert.False(t, m.RTryLockWithTimeout(5*time.Millisecond))

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.True(t, m.RTryLockWithTimeout(15*time.Millisecond))
	m.RUnlock()

	m.Lock()
	time.AfterFunc(10*time.Millisecond, func() {
		m.Unlock()
	})
	assert.False(t, m.RTryLockWithTimeout(5*time.Millisecond))

	assert.Panics(t, m.RUnlock)
}
