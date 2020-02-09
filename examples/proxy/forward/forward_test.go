package forward

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"testing"
)

func echoServer(network, addr string) {
	if len(network) == 0 {
		network = "tcp"
	}
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Println("server error:", err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go echo(conn)
		}
	}
}

func echo(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	for {
		buf := make([]byte, 64)
		rLen, err := conn.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("Read error:", err)
			continue
		}

		_, err = conn.Write(buf[:rLen])
		if err != nil {
			log.Println("Write error:", err)
			continue
		}
	}
}

func TestNew(t *testing.T) {
	_, err := New("listen", "forward", "network")
	assert.Error(t, err)
}

func TestProxy_Run(t *testing.T) {
	listen := "127.0.0.1:6666"
	forward := "127.0.0.1:8888"
	network := ""
	go echoServer(network, forward)
	proxy, err := New(listen, forward, network)
	assert.NoError(t, err)

	go proxy.Run()

	conn, err := net.Dial("tcp", listen)
	assert.NoError(t, err)
	_, err = conn.Write([]byte("hello"))
	assert.NoError(t, err)
	buf := make([]byte, 10)
	rLen, err := conn.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 5, rLen)
	assert.Equal(t, "hello", string(buf[:rLen]))
	err = conn.Close()
	assert.NoError(t, err)
}
