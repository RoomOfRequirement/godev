package limiter

import (
	"math"
	"testing"
	"time"
)

func TestLimiterBasic(t *testing.T) {
	f := func(rps, burst int) {
		defer func() {
			if r := recover(); r != nil {
				if r.(string) != "rps and burst should > 0" {
					t.Fatal(r)
				}
			}
		}()

		_, _ = Limiter(rps, burst)
	}
	f(-1, -1)
	f(-1, 1)
	f(1, -1)
}

func TestLimiter(t *testing.T) {
	tokenBucket, stopChan := Limiter(100, 10)
	defer func() {
		stopChan <- true
	}()
	cnt := -10
	ticker := time.NewTicker(1 * time.Second)

Loop:
	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			break Loop
		case <-tokenBucket:
			cnt++
		}
	}
	if !almostEqual(cnt, 100, 5) {
		t.Fatal(cnt)
	}
}

func TestLimiterStop(t *testing.T) {
	tokenBucket, stopChan := Limiter(100, 10)
	cnt := -10
	ticker1 := time.NewTicker(20 * time.Millisecond)
	ticker2 := time.NewTicker(1 * time.Second)
Loop:
	for {
		select {
		case <-ticker1.C:
			ticker1.Stop()
			stopChan <- true
		case <-ticker2.C:
			ticker2.Stop()
			break Loop
		case <-tokenBucket:
			cnt++
		}
	}
	// GHz = 10 ^ 9 Hz
	if !farLarger(cnt, 100, 9) {
		t.Fatal(cnt)
	}
}

func almostEqual(a, b, tolerance int) bool {
	if a-b < tolerance && b-a < tolerance {
		return true
	}
	return false
}

func farLarger(a, b, exp int) bool {
	if a - b > int(math.Pow(10., float64(exp))) || b - a < int(math.Pow(10., float64(exp))) {
		return true
	}
	return false
}
