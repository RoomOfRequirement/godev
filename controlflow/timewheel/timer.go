package timewheel

import "time"

// Timer ...
//	accuracy: millisecond
type Timer struct {
	task       *Task
	rounds     int
	executedAt int64 // milliseconds from (1970-1-1, UTC)
	interval   int64
}

// NowMS returns current time in milliseconds from (1970-1-1, UTC)
func NowMS() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// TimeMS returns time in milliseconds from (1970-1-1, UTC)
func TimeMS(t time.Time) int64 {
	return t.UTC().UnixNano() / 1e6
}

// DurationMS returns duration in milliseconds
func DurationMS(d time.Duration) int64 {
	return d.Milliseconds()
}

// NewTimer ...
func NewTimer(executedAt time.Time, fn func(...interface{}), args []interface{}) *Timer {
	return &Timer{
		task:       NewTask(fn, args),
		rounds:     0,
		executedAt: TimeMS(executedAt),
		interval:   0,
	}
}

// NewTimerWithRepeat ...
func NewTimerWithRepeat(executedAt time.Time, interval time.Duration, repeat int, fn func(...interface{}), args []interface{}) *Timer {
	if repeat < 1 {
		return NewTimer(executedAt, fn, args)
	}
	return &Timer{
		task:       NewTask(fn, args),
		rounds:     repeat - 1,
		executedAt: TimeMS(executedAt),
		interval:   DurationMS(interval),
	}
}

// Run ...
func (t *Timer) Run() {
	go t.run()
}

func (t *Timer) run() {
	now := NowMS()
	if t.executedAt > now {
		time.Sleep(time.Duration(t.executedAt-now) * time.Millisecond)
	}
	t.task.Run()
	if t.rounds > 0 {
		t.executedAt += t.interval
		t.rounds--
		t.run()
	}
}
