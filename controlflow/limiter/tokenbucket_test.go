package limiter

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	tb, err := NewBucket(-1, 0)
	assert.Nil(t, tb)
	assert.Error(t, err)

	tb, err = NewBucket(10, 0.01)
	assert.NotNil(t, tb)
	assert.NoError(t, err)

	success, err := tb.Acquire(-1)
	assert.Error(t, err)
	assert.False(t, success)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		success, err := tb.Acquire(5)
		assert.NoError(t, err)
		assert.True(t, success)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		success, err := tb.Acquire(6)
		assert.NoError(t, err)
		assert.False(t, success)
	}()
	wg.Wait()
}
