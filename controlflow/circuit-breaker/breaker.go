package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type breaker struct {
	maxRetries           int
	resetCounterInterval time.Duration // for closed state
	openInterval         time.Duration // for open state
	tripWhen             func(c Counter) bool
	onStateChange        func(from, to State)

	sync.Mutex
	state     State
	counter   Counter
	nextReset time.Time // reset depends on state (closed/open)
}

// New creates new breaker from input settings
func New(settings Settings) Interface {
	return &breaker{
		maxRetries:           settings.MaxRetries,
		resetCounterInterval: settings.ResetCounterInterval,
		openInterval:         settings.OpenInterval,
		tripWhen:             settings.TripWhen,
		onStateChange:        settings.OnStateChange,
		state:                Closed,
		counter:              Counter{},
		nextReset:            time.Now().Add(settings.ResetCounterInterval),
	}
}

// State returns current state of breaker
func (b *breaker) State() State {
	b.Lock()
	defer b.Unlock()
	return b.getState(time.Now())
}

func (b *breaker) getState(now time.Time) State {
	switch b.state {
	case Closed:
		if now.After(b.nextReset) {
			b.resetCounter(now)
		}
	case Open: // open -> half-open
		if now.After(b.nextReset) {
			b.setState(HalfOpen, now)
		}
	}
	return b.state
}

func (b *breaker) setState(s State, now time.Time) {
	if s == b.state {
		return
	}
	from, to := b.state, s
	b.state = s
	b.resetCounter(now)
	b.onStateChange(from, to)
}

func (b *breaker) resetCounter(now time.Time) {
	b.counter.reset()
	switch b.state {
	case Closed:
		b.nextReset = now.Add(b.resetCounterInterval)
	case Open:
		b.nextReset = now.Add(b.openInterval)
	case HalfOpen:
		b.nextReset = time.Time{}
	}
}

// Trip trips breaker
func (b *breaker) Trip() {
	b.Lock()
	defer b.Unlock()
	b.setState(Open, time.Now())
}

// Reset resets breaker
func (b *breaker) Reset() {
	*b = breaker{
		maxRetries:           b.maxRetries,
		resetCounterInterval: b.resetCounterInterval,
		openInterval:         b.openInterval,
		tripWhen:             b.tripWhen,
		onStateChange:        b.onStateChange,
		state:                Closed,
		counter:              Counter{},
		nextReset:            time.Now().Add(b.resetCounterInterval),
	}
}

// Execute executes task func depends on breaker's state
func (b *breaker) Execute(taskFunc func() (interface{}, error)) (interface{}, error) {
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	state := b.getState(now)
	switch state {
	case Open:
		return nil, errors.New("circuit breaker is OPEN now")
	case HalfOpen:
		if b.counter.Tasks >= b.maxRetries {
			return nil, errors.New("circuit breaker is HALF-OPEN and reaches max retries")
		}
	}
	res, err := taskFunc()
	if err != nil {
		b.onFailure(state, now)
		return nil, err
	}
	b.onSuccess(state, now)
	return res, nil
}

func (b *breaker) onSuccess(state State, now time.Time) {
	b.counter.onSuccess()
	if state == HalfOpen && b.counter.ConsecutiveSuccesses >= b.maxRetries {
		b.setState(Closed, now)
	}
}

func (b *breaker) onFailure(state State, now time.Time) {
	b.counter.onFailure()
	if state == HalfOpen {
		b.setState(Open, now)
	} else if state == Closed && b.tripWhen(b.counter) {
		b.setState(Open, now)
	}
}
