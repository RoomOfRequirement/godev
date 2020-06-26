package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"
)

// ErrWaitTimeout ...
var ErrWaitTimeout = errors.New("timeout")

// Server ...
type Server struct {
	cm     *ConnectionManager
	socket *net.TCPListener
}

// New ...
func New(port int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("fail to resolve addr: %v", err)
	}
	sock, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("fail to listen tcp: %v", err)
	}
	return &Server{
		cm:     NewConnectionManager(),
		socket: sock,
	}, nil
}

// NewFromFD ...
func NewFromFD(fd uintptr) (*Server, error) {
	file := os.NewFile(fd, path.Join(os.TempDir(), "/sock-go-graceful-restart"))
	listener, err := net.FileListener(file)
	if err != nil {
		return nil, fmt.Errorf("fail to recover socket from file descriptor: %v", err.Error())
	}
	sock, ok := listener.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("file descriptor %d is not a valid TCP socket", fd)
	}
	return &Server{
		cm:     NewConnectionManager(),
		socket: sock,
	}, nil
}

// Stop ...
func (s *Server) Stop() {
	// Accept will instantly return a timeout error
	_ = s.socket.SetDeadline(time.Now())
}

// SocketFD ...
func (s *Server) SocketFD() (uintptr, error) {
	file, err := s.socket.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

// Wait ...
func (s *Server) Wait() {
	s.cm.Wait()
}

// WaitWithTimeout ...
func (s *Server) WaitWithTimeout(duration time.Duration) error {
	timeout := time.NewTimer(duration)
	wait := make(chan struct{})
	go func() {
		s.Wait()
		wait <- struct{}{}
	}()
	select {
	case <-timeout.C:
		return ErrWaitTimeout
	case <-wait:
		return nil
	}
}

// StartAcceptLoop ...
func (s *Server) StartAcceptLoop() {
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				log.Println("Graceful Server - [Info] stop accepting connections")
				return
			}
			log.Println("Graceful Server - [Error] fail to accept:", err)
		}
		go func() {
			s.cm.Add(1)
			s.handleConn(conn)
			s.cm.Done()
		}()
	}
}

// ping - pong
func (s *Server) handleConn(conn net.Conn) {
	tick := time.NewTicker(time.Second)
	buffer := make([]byte, 64)
	for {
		select {
		case <-tick.C:
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				log.Println("Graceful Server - [Error] fail to ping:", err)
				_ = conn.Close()
				return
			}
			log.Println("Graceful Server - [Info] ping")

			n, err := conn.Read(buffer)
			if err != nil {
				log.Println("Graceful Server - [Error] fail to read from socket:", err)
				_ = conn.Close()
				return
			}
			log.Printf("Graceful Server - [Info] OK: read %d bytes: '%s'\n", n, string(buffer[:n]))
		}
	}
}

// Addr ...
func (s *Server) Addr() net.Addr {
	return s.socket.Addr()
}

// ConnectionsCounter ...
func (s *Server) ConnectionsCounter() int {
	return s.cm.Counter
}
