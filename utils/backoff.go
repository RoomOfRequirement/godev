package utils

import "time"

// BackOff ...
// Simple back-off wrapper, complex one see `retry.go` in `experiment` package
func BackOff(initialBackOff, maxBackOff time.Duration, maxCalls int, f func() error) error {
	backOff := initialBackOff
	// maxBackOff at least >= initialBackOff
	if maxBackOff < initialBackOff {
		maxBackOff = initialBackOff
	}
	calls := 0
	for {
		err := f()
		// success
		if err == nil {
			return nil
		}
		calls++
		// reach max calls
		// if calls > maxCalls && maxCalls != 0 -> if maxCalls == 0 means call f until success
		// better not use it here
		// now if maxCalls == 0 means call f only once
		if calls > maxCalls {
			return err
		}
		// reach max backOff interval
		if backOff > maxBackOff {
			backOff = maxBackOff
		} else {
			// exponentially increase
			backOff *= 2
		}
		time.Sleep(backOff)
		// log.Printf("[BackOff %v] Retry after %v due to the Error: %v\n", calls, backOff, err)
	}
}
