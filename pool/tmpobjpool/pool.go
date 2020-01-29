package tmpobjpool

import (
	"sync"
)

// Pool struct based on sync.Pool
type Pool struct {
	pool map[int]*sync.Pool
}

// New creates a new pool
func New() *Pool {
	return &Pool{
		pool: make(map[int]*sync.Pool),
	}
}

// AddSize adds a empty obj container with input size
func (p *Pool) AddSize(size int) {
	p.pool[size] = new(sync.Pool)
}

// Get gets a obj container with input size
func (p *Pool) Get(size int) interface{} {
	if pool := p.pool[size]; pool != nil {
		return pool.Get()
	}
	return nil
}

// Put puts back x into its object container according to its size
// if no object container found, then create one and put x into it
func (p *Pool) Put(x interface{}, size int) {
	if pool := p.pool[size]; pool != nil {
		pool.Put(x)
	} else {
		p.pool[size] = new(sync.Pool)
		p.pool[size].Put(x)
	}
}
