package roundrobin

import "sync"

// Ngx struct
//	from: http://hg.nginx.org/nginx/rev/d05ab8793a69
type Ngx struct {
	sync.Mutex
	nodes []*node
	next  int
}

type node struct {
	address         string
	weight          int
	currentWeight   int
	effectiveWeight int
}

// NewNgx creates a new Ngx
func NewNgx() *Ngx {
	return &Ngx{
		Mutex: sync.Mutex{},
		nodes: nil,
		next:  0,
	}
}

// SetNodes sets all nodes
func (nx *Ngx) SetNodes(nodes []string, weights []int) error {
	if len(nodes) == 0 || len(nodes) != len(weights) {
		return ErrInvalidNodes
	}
	nx.Lock()
	defer nx.Unlock()
	nx.nodes = make([]*node, len(nodes))
	for i := 0; i < len(nodes); i++ {
		nx.nodes[i] = &node{
			address:         nodes[i],
			weight:          weights[i],
			currentWeight:   0,
			effectiveWeight: weights[i],
		}
	}
	return nil
}

// Next returns a selected node string
func (nx *Ngx) Next() string {
	nx.Lock()
	defer nx.Unlock()
	total := 0
	selected := (*node)(nil)
	for i := 0; i < len(nx.nodes); i++ {
		node := nx.nodes[i]
		if node == nil || node.weight == 0 {
			continue
		}
		node.currentWeight += node.effectiveWeight
		total += node.effectiveWeight
		// if update node, useless here, so comment it
		/*
			if node.effectiveWeight < node.weight {
				node.effectiveWeight++
			}
		*/
		if selected == nil || node.currentWeight > selected.currentWeight {
			selected = node
		}
	}
	if selected == nil {
		return ""
	}
	selected.currentWeight -= total
	return selected.address
}

// Reset rests next selection to beginning state
func (nx *Ngx) Reset() {
	nx.Lock()
	defer nx.Unlock()
	for _, node := range nx.nodes {
		node.effectiveWeight = node.weight
		node.currentWeight = 0
	}
}
