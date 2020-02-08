package roundrobin

import "errors"

// RoundRobin ...
type RoundRobin interface {
	// SetNodes sets all nodes
	SetNodes(nodes []string, weights []int) error
	// Next returns selected node
	Next() string
	// Reset rests next selection to beginning state
	Reset()
}

// ALGORITHM ...
type ALGORITHM int

const (
	// Simple means simple round robin
	Simple ALGORITHM = iota
	// LVS means weighted round robin from LVS
	LVS
	// Nginx means weighted round robin from Nginx
	Nginx
)

// ErrInvalidNodes ...
var ErrInvalidNodes = errors.New("invalid number of nodes")

// ErrUnsupportedAlgorithm ...
var ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")

// New returns a new RoundRobin
func New(algorithm ALGORITHM) (RoundRobin, error) {
	switch algorithm {
	case Simple:
		return NewSRR(), nil
	case LVS:
		return NewLVS(), nil
	case Nginx:
		return NewNgx(), nil
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}
