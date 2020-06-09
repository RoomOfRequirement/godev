package p2c

import (
	"errors"
	"github.com/spaolacci/murmur3"
	"hash/fnv"
	"math/rand"
	"sync"
	"time"
)

// ErrNoNodes ...
var ErrNoNodes = errors.New("no nodes inside")

// ErrNodeNotExist ...
var ErrNodeNotExist = errors.New("node not exist")

// P2C means power of two choices
type P2C struct {
	sync.Mutex
	rnd      *rand.Rand
	nodesSet map[string]uint64 // node: addr/load, load should be updated in real time in practice
	nodes    []string          // addr
}

// NewP2C ...
func NewP2C() *P2C {
	return &P2C{
		Mutex:    sync.Mutex{},
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
		nodesSet: make(map[string]uint64),
		nodes:    make([]string, 0),
	}
}

// AddNode adds one node with load
func (p *P2C) AddNode(addr string, load uint64) {
	p.Lock()
	defer p.Unlock()

	if _, found := p.nodesSet[addr]; found {
		return
	}
	p.nodesSet[addr] = load
	p.nodes = append(p.nodes, addr)
	return
}

// DeleteNode deletes one node
func (p *P2C) DeleteNode(addr string) {
	p.Lock()
	defer p.Unlock()

	if _, found := p.nodesSet[addr]; !found {
		return
	}
	delete(p.nodesSet, addr)
	for i := range p.nodes {
		if p.nodes[i] == addr {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			return
		}
	}
}

func (p *P2C) get2(key string) (string, string) {
	l := len(p.nodes)
	if len(key) == 0 {
		return p.nodes[p.rnd.Intn(l)], p.nodes[p.rnd.Intn(l)]
	}
	data := []byte(key)
	h := fnv.New32()
	_, _ = h.Write(data)
	m := murmur3.New32()
	_, _ = m.Write(data)
	return p.nodes[int(h.Sum32())%l], p.nodes[int(m.Sum32())%l]
}

// Get gets one node with less load for input key based on two hash
//	if no key, then randomly picks two
func (p *P2C) Get(key string) (addr string, err error) {
	p.Lock()
	defer p.Unlock()

	if len(p.nodes) == 0 {
		return "", ErrNoNodes
	}
	n1, n2 := p.get2(key)
	// choose the one with less load
	if p.nodesSet[n1] <= p.nodesSet[n2] {
		p.nodesSet[n1]++
		return n1, nil
	}
	p.nodesSet[n2]++
	return n2, nil
}

// IncrLoad increases one load for node with addr
func (p *P2C) IncrLoad(addr string) error {
	p.Lock()
	defer p.Unlock()

	if _, found := p.nodesSet[addr]; found {
		p.nodesSet[addr]++
		return nil
	}
	return ErrNodeNotExist
}

// DecrLoad decreases one load for node with addr
func (p *P2C) DecrLoad(addr string) error {
	p.Lock()
	defer p.Unlock()

	if load, found := p.nodesSet[addr]; found {
		if load > 0 {
			p.nodesSet[addr]--
		}
		return nil
	}
	return ErrNodeNotExist
}

// UpdateLoad updates load for node with addr
func (p *P2C) UpdateLoad(addr string, load uint64) error {
	p.Lock()
	defer p.Unlock()

	if _, found := p.nodesSet[addr]; found {
		p.nodesSet[addr] = load
		return nil
	}
	return ErrNodeNotExist
}

// GetLoad gets load for node with addr
func (p *P2C) GetLoad(addr string) (uint64, error) {
	p.Lock()
	defer p.Unlock()

	if load, found := p.nodesSet[addr]; found {
		return load, nil
	}
	return 0, ErrNodeNotExist
}
