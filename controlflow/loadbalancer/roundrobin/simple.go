package roundrobin

import (
	"sync"
)

// SimpleRoundRobin struct
type SimpleRoundRobin struct {
	sync.Mutex
	nodes []string
	next  int
}

// NewSRR creates a new SimpleRoundRobin
func NewSRR() *SimpleRoundRobin {
	return &SimpleRoundRobin{
		Mutex: sync.Mutex{},
		nodes: nil,
		next:  0,
	}
}

// SetNodes sets all nodes
func (srr *SimpleRoundRobin) SetNodes(nodes []string, weights []int) error {
	if len(nodes) == 0 {
		return ErrInvalidNodes
	}
	srr.Lock()
	defer srr.Unlock()
	srr.nodes = nodes
	return nil
}

// Next returns a selected node string
func (srr *SimpleRoundRobin) Next() string {
	srr.Lock()
	defer srr.Unlock()
	selected := srr.nodes[srr.next]
	srr.next = (srr.next + 1) % len(srr.nodes)
	return selected
}

// Reset rests next selection to beginning state
func (srr *SimpleRoundRobin) Reset() {
	srr.Lock()
	defer srr.Unlock()
	srr.next = 0
}
