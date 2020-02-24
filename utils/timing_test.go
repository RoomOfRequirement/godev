package utils

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	f := func() {
		defer RoughTiming(time.Now(), "test")
		// time.Sleep(100 * time.Millisecond)
		return
	}
	f()
}
