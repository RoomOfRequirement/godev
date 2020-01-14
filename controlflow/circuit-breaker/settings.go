package circuitbreaker

import (
	"fmt"
	"log"
	"time"
)

// Settings of breaker
type Settings struct {
	MaxRetries           int
	ResetCounterInterval time.Duration
	OpenInterval         time.Duration
	TripWhen             func(c Counter) bool
	OnStateChange        func(from, to State)
}

// NewSettings creates new settings of breaker
func NewSettings(maxRetries int, resetCounterInterval, openInterval time.Duration, tripWhen func(c Counter) bool, onStateChange func(from, to State)) Settings {
	if maxRetries <= 0 {
		maxRetries = 1
	}
	if resetCounterInterval == 0 {
		resetCounterInterval = 60 * time.Second
	}
	if openInterval == 0 {
		openInterval = 60 * time.Second
	}
	if tripWhen == nil {
		tripWhen = func(c Counter) bool {
			if 2*c.TotalFailures > c.Tasks || c.ConsecutiveFailures > 10 {
				return true
			}
			return false
		}
	}
	if onStateChange == nil {
		onStateChange = func(from, to State) {
			log.Println("Circuit Breaker changes state from", fmt.Sprintf("%s", from), "to", fmt.Sprintf("%s", to))
		}
	}
	return Settings{
		MaxRetries:           maxRetries,
		ResetCounterInterval: resetCounterInterval,
		OpenInterval:         openInterval,
		TripWhen:             tripWhen,
		OnStateChange:        onStateChange,
	}
}

// Counter of breaker
type Counter struct {
	Tasks                int
	TotalSuccesses       int
	TotalFailures        int
	ConsecutiveSuccesses int
	ConsecutiveFailures  int
}

func (c *Counter) onSuccess() {
	c.Tasks++
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counter) onFailure() {
	c.Tasks++
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counter) reset() {
	c.Tasks = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}
