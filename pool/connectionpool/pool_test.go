package connectionpool

import (
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

var serverOnce sync.Once
var serverAddr net.Addr

var factory = func() (net.Conn, error) {
	serverOnce.Do(func() {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		serverAddr = l.Addr()

		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}

				go func() {
					_, _ = io.Copy(os.Stdout, conn)
				}()
			}
		}()
	})

	return net.Dial(serverAddr.Network(), serverAddr.String())
}

func TestNew(t *testing.T) {
	_, err := New(-10, factory)
	if err != ErrInvalidLimit {
		t.Fatalf("New error: %s", err)
	}

	_, err = New(10, nil)
	if err != ErrInvalidFactory {
		t.Fatalf("New error: %s", err)
	}

	_, err = New(10, factory)
	if err != nil {
		t.Fatalf("New error: %s", err)
	}
}

func TestPool_Get(t *testing.T) {
	p, err := New(10, factory)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	// get with invalid factory
	err = p.ChangeFactory(func() (conn net.Conn, err error) {
		return nil, ErrInvalidFactory
	})
	if err != nil {
		t.Fatal(err)
	}
	conn, err := p.Get()
	if conn != nil || err != ErrInvalidFactory {
		t.Fatal(conn, err)
	}
	// after test, change factory back
	err = p.ChangeFactory(factory)
	if err != nil {
		t.Fatal(err)
	}

	// get one
	_, err = p.Get()
	if err != nil {
		t.Fatal(err)
	}

	// get all
	var wg sync.WaitGroup
	for i := 0; i < 10-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Get()
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()

	// test empty
	c, err := p.Get()
	if c != nil || err != nil {
		t.Fatal(c, err)
	}
	if len(p.semaphore) != 10 || len(p.idleConnections) != 0 {
		t.Fatal(len(p.semaphore), len(p.idleConnections))
	}

	// test closed
	p.Close()
	c, err = p.Get()
	if c != nil || err != ErrClosed {
		t.Fatal(c, err)
	}
}

func TestPool_Put(t *testing.T) {
	p, err := New(10, factory)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	// create and get all
	conns := make([]net.Conn, 10)
	for i := range conns {
		conn, _ := p.Get()
		conns[i] = conn
	}
	if len(p.semaphore) != 10 || len(p.idleConnections) != 0 {
		t.Fatal(len(p.semaphore), len(p.idleConnections))
	}

	// put all
	for _, conn := range conns {
		if err := p.Put(conn); err != nil {
			t.Fatal(err)
		}
	}
	if len(p.semaphore) != 0 || len(p.idleConnections) != 10 {
		t.Fatal(len(p.semaphore), len(p.idleConnections))
	}

	// put one more
	conn, err := factory()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Put(conn); err != nil {
		t.Fatal(err)
	}
	_, err = conn.Read(make([]byte, 0))
	if err == nil {
		t.Fatal("not close conn when pool full")
	} else {
		t.Logf("conn closed: %v", err)
	}

	// get one
	conn, err = p.Get()
	if err != nil {
		t.Fatal(err)
	}
	// put invalid conn
	if err := p.Put(nil); err != ErrInvalidConnection {
		t.Fatal(err)
	}

	// put on closed
	p.Close()
	conn, err = factory()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Put(conn); err != ErrClosed {
		t.Fatal(err)
	}
	_, err = conn.Read(make([]byte, 0))
	if err == nil {
		t.Fatal("not close conn when pool full")
	} else {
		t.Logf("conn closed: %v", err)
	}
}

func TestPool_ChangeFactory(t *testing.T) {
	p, err := New(10, factory)
	if err != nil {
		t.Fatal(err)
	}
	err = p.ChangeFactory(nil)
	if err != ErrInvalidFactory {
		t.Fatal(err)
	}
	err = p.ChangeFactory(factory)
	if err != nil {
		t.Fatal(err)
	}

	p.Close()
	err = p.ChangeFactory(factory)
	if err != ErrClosed {
		t.Fatal(err)
	}
}

func TestPoolConcurrent_Get(t *testing.T) {
	p, _ := New(20, factory)
	defer p.Close()
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				_, _ = p.Get()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				wg.Done()
			}()
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			_, _ = p.Get()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestPoolConcurrent_GetPut(t *testing.T) {
	p, _ := New(20, factory)
	defer p.Close()
	conn := make(chan net.Conn)

	go func() {
		p.Close()
	}()

	for i := 0; i < 20; i++ {
		go func() {
			c, _ := p.Get()
			conn <- c
		}()

		go func() {
			c := <- conn
			if c == nil {
				return
			}
			_ = p.Put(c)
		}()
	}
}
