package loadbalancer

import (
	"errors"
	"sync"
)

// ErrInvalidNodes ...
var ErrInvalidNodes = errors.New("invalid number of nodes")

// RoundRobin struct
type RoundRobin struct {
	sync.Mutex
	nodes []string
	next  int
}

// NewRR creates a new RoundRobin
func NewRR(nodes []string) (*RoundRobin, error) {
	if len(nodes) == 0 {
		return nil, ErrInvalidNodes
	}
	return &RoundRobin{
		Mutex: sync.Mutex{},
		nodes: nodes,
		next:  0,
	}, nil
}

// Next returns a selected node string
func (rr *RoundRobin) Next() string {
	rr.Lock()
	defer rr.Unlock()
	selected := rr.nodes[rr.next]
	rr.next = (rr.next + 1) % len(rr.nodes)
	return selected
}
