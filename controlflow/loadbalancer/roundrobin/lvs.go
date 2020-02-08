package roundrobin

import (
	"sync"
)

// LVSrr struct
//	from: http://www.linuxvirtualserver.org/docs/scheduling.html
type LVSrr struct {
	sync.Mutex
	nodes         []string
	weights       []int
	next          int
	currentWeight int
	gcd           int
}

// NewLVS creates a new LVSrr
func NewLVS() *LVSrr {
	return &LVSrr{
		Mutex:         sync.Mutex{},
		nodes:         nil,
		weights:       nil,
		next:          -1,
		currentWeight: 0,
		gcd:           0,
	}
}

// SetNodes sets all nodes
func (wrr *LVSrr) SetNodes(nodes []string, weights []int) error {
	if len(nodes) == 0 || len(nodes) != len(weights) {
		return ErrInvalidNodes
	}
	wrr.Lock()
	defer wrr.Unlock()
	wrr.nodes = nodes
	wrr.weights = weights
	wrr.gcd = gcdArray(weights)
	return nil
}

// Next returns a selected node string
func (wrr *LVSrr) Next() string {
	wrr.Lock()
	defer wrr.Unlock()
	next := wrr.next
	for {
		next = (next + 1) % len(wrr.nodes)
		if next == 0 {
			cw := wrr.currentWeight - wrr.gcd
			if cw <= 0 {
				cw = maxArray(wrr.weights)
				if cw == 0 {
					return ""
				}
			}
			wrr.currentWeight = cw
		}
		if wrr.weights[next] >= wrr.currentWeight {
			wrr.next = next
			break
		}
	}
	return wrr.nodes[next]
}

// Reset rests next selection to beginning state
func (wrr *LVSrr) Reset() {
	wrr.Lock()
	defer wrr.Unlock()
	wrr.next = -1
	wrr.currentWeight = 0
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
