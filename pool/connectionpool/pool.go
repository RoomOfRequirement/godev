package connectionpool

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// Pool struct without explict mutex lock
type Pool struct {
	semaphore       chan token
	idleConnections chan net.Conn

	factory chan Factory

	closeOnce sync.Once
	closed    int32
}

// Factory function for creating new connection
type Factory func() (net.Conn, error)

type token struct{}

// ErrClosed ...
var ErrClosed = errors.New("connection pool closed")

// ErrInvalidLimit ...
var ErrInvalidLimit = errors.New("invalid limit number")

// ErrInvalidFactory ...
var ErrInvalidFactory = errors.New("invalid factory function for creating new connection")

// ErrInvalidConnection ...
var ErrInvalidConnection = errors.New("invalid nil connection")

// New creates a new connection pool
func New(limit int, factory Factory) (*Pool, error) {
	if limit < 1 {
		return nil, ErrInvalidLimit
	}
	if factory == nil {
		return nil, ErrInvalidFactory
	}

	p := &Pool{
		semaphore:       make(chan token, limit),
		idleConnections: make(chan net.Conn, limit),
		factory:         make(chan Factory, 1),
		closeOnce:       sync.Once{},
		closed:          0,
	}
	p.factory <- factory
	return p, nil
}

// Get gets one conn from pool
func (p *Pool) Get() (net.Conn, error) {
	if p.Closed() {
		return nil, ErrClosed
	}
	select {
	// get one
	case conn := <-p.idleConnections:
		return conn, nil
	// not up to semaphore
	case p.semaphore <- token{}:
		f := <-p.factory
		conn, err := f()
		p.factory <- f
		if err != nil {
			<-p.semaphore
		}
		return conn, err
	// empty
	default:
		return nil, nil
	}
}

// Put puts a conn into the pool
//	if pool is closed or full, conn will be closed
func (p *Pool) Put(conn net.Conn) error {
	// pool is closed, close input conn and return ErrClosed
	if p.Closed() {
		defer func() {
			_ = conn.Close()
		}()
		return ErrClosed
	}
	if conn == nil {
		return ErrInvalidConnection
	}
	select {
	// pool not full
	case <-p.semaphore:
		p.idleConnections <- conn
		return nil
	// pool is full, close input conn
	default:
		return conn.Close()
	}
}

// Close closes the pool
func (p *Pool) Close() {
	p.closeOnce.Do(func() {
		atomic.StoreInt32(&p.closed, 1)
		conns := p.idleConnections
		p.idleConnections = nil
		close(p.semaphore)
		close(conns)
		for conn := range conns {
			_ = conn.Close()
		}
	})
}

// Closed returns true if this connection pool has been closed
func (p *Pool) Closed() bool {
	return atomic.LoadInt32(&p.closed) != 0
}

// ChangeFactory changes factory func
func (p *Pool) ChangeFactory(factory Factory) error {
	if factory == nil {
		return ErrInvalidFactory
	}
	if p.Closed() {
		return ErrClosed
	}
	<-p.factory
	p.factory <- factory
	return nil
}
