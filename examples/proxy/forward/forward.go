package forward

import (
	"godev/controlflow/pipe"
	"io"
	"log"
	"net"
	"sync"
)

// Proxy ...
type Proxy struct {
	listen   string
	forward  string
	network  string
	listener net.Listener
}

// New ...
func New(listenAddr, forwardAddr, network string) (*Proxy, error) {
	// default TCP
	if len(network) == 0 {
		network = "tcp"
	}
	l, err := net.Listen(network, listenAddr)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		listen:   listenAddr,
		forward:  forwardAddr,
		network:  network,
		listener: l,
	}, nil
}

// Run ...
func (p *Proxy) Run() {
	wg := sync.WaitGroup{}
	go p.run(wg)
	wg.Wait()
	return
}

func (p *Proxy) run(wg sync.WaitGroup) {
	for {
		comeFrom, err := p.listener.Accept()
		if err != nil {
			log.Printf("listener failed to accept from: %s, with error: %s\n", comeFrom.RemoteAddr(), err)
			return
		}
		log.Printf("accepted from: %s", comeFrom.RemoteAddr())
		wg.Add(1)
		go func(network, forward string) {
			defer wg.Done()
			defer func() {
				_ = comeFrom.Close()
			}()

			forwardTo, err := net.Dial(network, forward)
			if err != nil {
				log.Printf("connection failed from %s to %s: %s\n", comeFrom.RemoteAddr(), p.forward, err)
				return
			}
			defer func() {
				_ = forwardTo.Close()
			}()

			if err := pipe.Pipe(comeFrom, forwardTo); err != nil && err != io.ErrClosedPipe {
				log.Printf("pipe failed: %s from %s to %s\n", err, comeFrom.RemoteAddr(), forwardTo.RemoteAddr())
			} else {
				log.Printf("connection closed from %s to %s\n", comeFrom.RemoteAddr(), forwardTo.RemoteAddr())
			}
			return
		}(p.network, p.forward)
	}
}
