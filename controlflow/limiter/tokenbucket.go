package limiter

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket ...
type TokenBucket struct {
	sync.Mutex

	capacity      int
	tokens        int
	refillRate    float64 // per sec
	lastTimestamp int64
}

// NewBucket ...
func NewBucket(capacity int, refillRate float64) (*TokenBucket, error) {
	if capacity < 1 || refillRate < 0 {
		return nil, fmt.Errorf("invalid capacity: %d or refill rate %0.2f", capacity, refillRate)
	}
	return &TokenBucket{
		capacity:      capacity,
		tokens:        0,
		refillRate:    refillRate,
		lastTimestamp: 0,
	}, nil
}

// Acquire ...
func (tb *TokenBucket) Acquire(n int) (bool, error) {
	if n < 1 {
		return false, fmt.Errorf("invalid acquire amount: %d", n)
	}
	tb.Lock()
	defer tb.Unlock()
	now := time.Now().UTC().Unix() // s
	tb.tokens = tb.tokens + int(float64(now-tb.lastTimestamp)*tb.refillRate)
	tb.lastTimestamp = now
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	if tb.tokens > n {
		tb.tokens -= n
		return true, nil
	}
	return false, nil
}
