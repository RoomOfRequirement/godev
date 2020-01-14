package circuitbreaker

import "strconv"

// Interface ...
//	Circuit Breaker Pattern: https://docs.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
type Interface interface {
	State() State
	Trip()
	Reset()

	Execute(taskFunc func() (interface{}, error)) (interface{}, error)
}

// State indicates circuit breaker state
type State int

const (
	// Open state means service not available
	Open State = iota
	// HalfOpen state means testing on service availability
	HalfOpen
	// Closed state means service available
	Closed
)

// String returns state string for print
func (s State) String() string {
	switch s {
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF-OPEN"
	case Closed:
		return "CLOSED"
	default:
		return "Unknown State: " + strconv.Itoa(int(s))
	}
}
