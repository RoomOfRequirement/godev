package spinlock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	// TryLock
	l := &Lock{}
	assert.True(t, l.TryLock())

	// IsLocked
	assert.True(t, l.IsLocked())

	// TryUnlock
	assert.True(t, l.TryUnlock())

	// SpinLock
	b := l.SpinLock(time.Millisecond, 5*time.Millisecond)
	assert.True(t, l.IsLocked())
	assert.True(t, b)
	b = l.SpinLock(time.Millisecond, 5*time.Millisecond)
	assert.True(t, l.IsLocked())
	assert.False(t, b)

	// ForceUnlock
	l.ForceUnlock()
	assert.False(t, l.IsLocked())

	// ForceLock
	l.ForceLock()
	assert.True(t, l.IsLocked())
	time.AfterFunc(100*time.Millisecond, l.ForceUnlock)
	l.ForceLock()
}
