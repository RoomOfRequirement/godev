package timewheel

import (
	"sync"
	"time"
)

// TimeWheel ...
type TimeWheel struct {
	sync.Mutex

	name string // name it

	tick           int64               // one tick to move one slot (in millisecond)
	slots          int                 // slots in time wheel, e.g. 60 for second
	slotCap        int                 // max timer number inside one slot
	currentIdx     int                 // current slot idx
	slotWithTimers []map[uint64]*Timer // slot -> timer by unique id

	childWheel *TimeWheel // child wheel with smaller tick

	stop chan struct{}  // stop chan
	stopOnce sync.Once
}

// NewTimeWheel ...
func NewTimeWheel(name string, tick int64, slots int, slotCap int) *TimeWheel {
	tw := &TimeWheel{
		Mutex:        sync.Mutex{},
		name:           name,
		tick:           tick,
		slots:          slots,
		slotCap:        slotCap,
		currentIdx:     0,
		slotWithTimers: make([]map[uint64]*Timer, slots),
		childWheel:     nil,
		stop: make(chan struct{}, 1),
		stopOnce: sync.Once{},
	}
	for i := range tw.slotWithTimers {
		tw.slotWithTimers[i] = make(map[uint64]*Timer, slotCap)
	}
	return tw
}

// AddTimer ...
func (tw *TimeWheel) AddTimer(uid uint64, timer *Timer) error {
	tw.Lock()
	defer tw.Unlock()

	d := timer.executedAt - NowMS()
	if d >= tw.tick {
		idx := (int(d/tw.tick) + tw.currentIdx) % tw.slots
		tw.slotWithTimers[idx][uid] = timer
		return nil
	} else if d < tw.tick && tw.childWheel == nil {
		// tw is smaller time wheel, let d = tw.tick and add to next round
		// if add to current round, it may not able to run it in time
		tw.slotWithTimers[(tw.currentIdx+1)%tw.slots][uid] = timer
		return nil
	}
	// add to its child time wheel
	return tw.childWheel.AddTimer(uid, timer)
}

// RemoveTimer ...
func (tw *TimeWheel) RemoveTimer(uid uint64) {
	tw.Lock()
	defer tw.Unlock()

	for i := range tw.slotWithTimers {
		if _, found := tw.slotWithTimers[i][uid]; found {
			delete(tw.slotWithTimers[i], uid)
		}
	}
}

// SetChild ...
func (tw *TimeWheel) SetChild(child *TimeWheel) {
	tw.Lock()
	defer tw.Unlock()

	tw.childWheel = child
}

// Run ...
func (tw *TimeWheel) Run() {
	go func() {
		for {
			select {
			case <- tw.stop:
				return
			default:
				tw.Lock()
				currentTimers := tw.slotWithTimers[tw.currentIdx]
				for _, timer := range currentTimers {
					// execute
					timer.Run()
				}
				// clear
				tw.slotWithTimers[tw.currentIdx] = make(map[uint64]*Timer)
				// move forward
				tw.currentIdx++
				tw.Unlock()

				time.Sleep(time.Duration(tw.tick) * time.Millisecond)
			}
		}
	}()
}

// Stop ...
func (tw *TimeWheel) Stop() {
	tw.stopOnce.Do(func() {
		tw.stop <- struct{}{}
	})
}
