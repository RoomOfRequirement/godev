package flow

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

// WaitSig struct
type WaitSig struct {
	wg   *sync.WaitGroup
	sigs chan os.Signal
	done chan struct{}
}

// NewWait creates a new `WaitSig` with input os signals (e.g. SIGINT / SIGTERM)
//	when corresponding type signals are triggered, `WaitSig` will exit
func NewWait(sigs ...os.Signal) *WaitSig {
	w := &WaitSig{
		wg:   &sync.WaitGroup{},
		sigs: make(chan os.Signal, 1),
		done: make(chan struct{}),
	}
	signal.Notify(w.sigs, sigs...)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.sigs:
				return
			case <-w.done:
				return
			}
		}
	}()
	return w
}

// WaitFor waits for signal to call func
func (w *WaitSig) WaitFor(f func() error) error {
	w.wg.Wait()
	return f()
}

// WaitForOrAfter waits for signal to call func or auto call func after interval
func (w *WaitSig) WaitForOrAfter(f func() error, interval time.Duration) error {
	time.AfterFunc(interval, func() {
		w.done <- struct{}{}
	})
	w.wg.Wait()
	return f()
}
