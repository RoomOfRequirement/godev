package loadbalancer

import (
	"sync"
)

// RoundRobinWeighted struct
type RoundRobinWeighted struct {
	sync.Mutex
	nodes         []string
	weights       []int
	next          int
	currentWeight int
	gcd           int
}

// NewRRW creates a new RoundRobinWeighted
func NewRRW(nodes []string, weights []int) (*RoundRobinWeighted, error) {
	if len(nodes) == 0 || len(nodes) != len(weights) {
		return nil, ErrInvalidNodes
	}
	return &RoundRobinWeighted{
		Mutex:         sync.Mutex{},
		nodes:         nodes,
		weights:       weights,
		next:          -1,
		currentWeight: 0,
		gcd:           gcdArray(weights),
	}, nil
}

// Next returns a selected node string
func (rrw *RoundRobinWeighted) Next() string {
	rrw.Lock()
	defer rrw.Unlock()
	next := rrw.next
	for {
		next = (next + 1) % len(rrw.nodes)
		if next == 0 {
			cw := rrw.currentWeight - rrw.gcd
			if cw <= 0 {
				cw = maxArray(rrw.weights)
				if cw == 0 {
					return ""
				}
			}
			rrw.currentWeight = cw
		}
		if rrw.weights[next] >= rrw.currentWeight {
			rrw.next = next
			break
		}
	}
	return rrw.nodes[next]
}

// Euclidean
func gcd(a, b int) int {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}

// gcd(a,b,c)=gcd(gcd(a,b),c)
func gcdArray(arr []int) int {
	ret := arr[0]
	for i := 1; i < len(arr); i++ {
		ret = gcd(ret, arr[i])
		if ret == 1 {
			return ret
		}
	}
	return ret
}

func maxArray(arr []int) int {
	max := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}
