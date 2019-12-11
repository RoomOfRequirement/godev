package limiter

import (
	"fmt"
	"os"
	"time"
)

// Limiter creates a token bucket for limiting resource visiting rate and burst size
//	notice: if trigger stopChan, then it becomes unlimited (read on closed will NOT block and return `zero-value` if empty)
// it is only a simple buffered chan
// for more complex usage, try official rate package: https://godoc.org/golang.org/x/time/rate
func Limiter(rps, burst int) (tokenBucket <-chan struct{}, stopChan chan<- bool) {
	if rps < 1 || burst < 1 {
		panic("rps and burst should > 0")
	}
	tb := make(chan struct{}, burst)
	sc := make(chan bool, 1)
	// initially fill bucket
	for i := 0; i < burst; i++ {
		tb <- struct{}{}
	}

	// create a ticker with interval = 1 / rps
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	// refill bucket
	go func() {
		defer func() {
			close(tb)
			ticker.Stop()
		}()
		for range ticker.C {
			select {
			case tb <- struct{}{}:
				// refill
			case stop := <-sc:
				if stop {
					_, _ = fmt.Fprintln(os.Stderr, "limiter stopped")
					return
				}
			default:
				// not refill until next interval
			}
		}
	}()

	return tb, sc
}
