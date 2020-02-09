package pipe

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestPipe(t *testing.T) {
	server1, client1 := net.Pipe()
	server2, client2 := net.Pipe()
	defer func() {
		_ = server1.Close()
		_ = client1.Close()
		_ = server2.Close()
		_ = client2.Close()
		// time.Sleep(20 * time.Millisecond)
		// to see "io: read/write on closed pipe"
	}()
	go func() {
		err := Pipe(server1, server2)
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err := client1.Write([]byte("hello"))
	assert.NoError(t, err)
	buf := make([]byte, 10)
	l, err := client2.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(buf[:l]))

	_, err = client2.Write([]byte("world"))
	assert.NoError(t, err)
	buf = make([]byte, 10)
	l, err = client1.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, "world", string(buf[:l]))
}
