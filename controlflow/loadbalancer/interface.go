package loadbalancer

// LoadBalancer interface
type LoadBalancer interface {
	Select(nodes []string, key string) (selectedNode string, err error)
}

type loadBalancer struct {
	nodes []string
	ch    *ConsistentHash
}

func (lb *loadBalancer) Select(nodes []string, key string) (selectedNode string, err error) {
	lb.nodes = nodes
	lb.ch.Set(lb.nodes)
	return lb.ch.Get(key)
}

// NewWithConsistentHash returns a LoadBalancer based on ConsistentHash
func NewWithConsistentHash() LoadBalancer {
	return &loadBalancer{
		nodes: nil,
		ch:    NewCH(nil),
	}
}
