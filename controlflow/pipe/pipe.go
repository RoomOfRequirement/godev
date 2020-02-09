package pipe

import (
	"io"
	"net"
)

// Pipe bi-directionally copies data between two connections, it returns io.ErrClosedPipe when connections close
//	here simply uses io.Copy, whose buffer gc after one request
//	if you really cares about allocations, you should use buffer pool
//	like this: https://github.com/golang/go/commit/492a62e945555bbf94a6f9dd6d430f712738c5e0
//	for in-memory cases, you may want these two pipes: net.Pipe() (bi-directional) or io.Pipe() (single directional)
func Pipe(conn1 net.Conn, conn2 net.Conn) error {
	errChan := make(chan error, 1)
	dataCopy := func(w, r net.Conn) {
		_, err := io.Copy(w, r)
		errChan <- err
	}
	go dataCopy(conn1, conn2)
	go dataCopy(conn2, conn1)
	return <-errChan
}
