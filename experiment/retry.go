package experiment

import (
	"log"
	"math"
	"sync"
	"time"
)

// RunWithRetry ...
func RunWithRetry(fn func() (success bool, err error), notifyFn func(trialNo int, err error), maxTrials int) {
	var (
		success bool
		err     error
	)
	for i := 0; i < maxTrials; i++ {
		// succeed
		if success, err = fn(); success {
			return
		}
		// fail
		if notifyFn != nil {
			notifyFn(i, err)
		}
	}
}

// RunWithBackOff ...
//	retry trials -> wait for some time and adjust its trial interval -> retry -> ... -> reset trial interval -> retry -> ...
func RunWithBackOff(fn func() (success bool, err error), notifyFn func(trialNo int, wait time.Duration, err error), trials int, minInterval, maxInterval, wait, reset time.Duration) (cancel func()) {
	off := &backOff{
		mu:          sync.RWMutex{},
		trials:      trials,
		min:         float64(minInterval),
		max:         float64(maxInterval),
		adjustedMin: float64(minInterval),
		tuneFactor:  0,
		round:       0,
		accumulated: 0,
		wg:          sync.WaitGroup{},
		done:        make(chan struct{}),
	}
	off.tuneFactor = math.Pow(off.max/off.min, 1/float64(off.trials-1))
	off.wg.Add(1)
	go off.auto(wait, reset)
	go off.run(fn, notifyFn)
	return func() {
		close(off.done)
	}
}

type backOff struct {
	mu sync.RWMutex

	trials                int
	min, max, adjustedMin float64
	tuneFactor            float64

	// auto-tune, changes trial rate according to former fail accumulated rate
	round       int
	accumulated float64
	wg          sync.WaitGroup
	done        chan struct{}
}

func (off *backOff) auto(wait, reset time.Duration) {
	defer off.wg.Done()
	for {
		w := time.NewTimer(wait)
		select {
		case <-off.done:
			w.Stop()
			return
		case <-w.C:
			// wait
		}
		off.mu.Lock()
		if off.round > 0 {
			off.adjustedMin = off.accumulated / float64(off.round)
			off.tuneFactor = math.Pow(off.max/off.adjustedMin, 1/float64(off.trials-1))
		}
		off.mu.Unlock()
		log.Println("tune factor")

		// reset timer to reset
		w.Reset(reset)
		select {
		case <-off.done:
			w.Stop()
			return
		case <-w.C:
			// reset
		}
		off.mu.Lock()
		off.adjustedMin = off.min
		off.tuneFactor = math.Pow(off.max/off.min, 1/float64(off.trials-1))
		off.mu.Unlock()
		log.Println("reset")
	}
}

func (off *backOff) run(fn func() (success bool, err error), notifyFn func(trialNo int, wait time.Duration, err error)) {
	defer func() {
		if _, open := <-off.done; open {
			close(off.done)
		}
		off.wg.Wait()
	}()
	notify := func(trialNo int, err error) {
		off.mu.Lock()
		off.round++
		wait := off.adjustedMin * math.Pow(off.tuneFactor, float64(trialNo-1))
		off.accumulated += wait
		off.mu.Unlock()
		if notifyFn != nil {
			notifyFn(trialNo, time.Duration(wait), err)
		}
		time.Sleep(time.Duration(wait))
	}
	RunWithRetry(fn, notify, off.trials)
}
