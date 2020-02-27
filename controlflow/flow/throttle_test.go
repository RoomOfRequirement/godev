package flow

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewThrottle(t *testing.T) {
	t.Parallel()

	t.Run("leading false", func(t *testing.T) {
		t.Parallel()

		testThrottle(t, false)
	})

	t.Run("leading true", func(t *testing.T) {
		t.Parallel()

		testThrottle(t, true)
	})
}

func testThrottle(t *testing.T, leading bool) {
	var counter uint64
	f := func() {
		atomic.AddUint64(&counter, 1)
	}
	throttled, _ := NewThrottle(context.TODO(), 100*time.Millisecond, leading)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			throttled(f)
		}
		time.Sleep(60 * time.Millisecond)
	}
	// let f executed
	time.Sleep(100 * time.Millisecond)
	c := int(atomic.LoadUint64(&counter))
	if c != 2 {
		t.Error("Expected count 2, was", c)
	}

	// cancel called, so f not executed
	throttled, cancel := NewThrottle(context.TODO(), 100*time.Millisecond, leading)
	cancel()
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			throttled(f)
		}
		time.Sleep(60 * time.Millisecond)
	}
	// let f executed
	time.Sleep(100 * time.Millisecond)
	c = int(atomic.LoadUint64(&counter))
	if c != 2 {
		t.Error("Expected count 2, was", c)
	}
}
