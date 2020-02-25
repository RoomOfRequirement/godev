package service

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// RunService ...
func RunService(name, addr, master string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("service listen at: %s\n", l.Addr().String())
	hb := NewHeartBeat(name, l.Addr().(*net.TCPAddr).Port, master)
	hb.Start()
	// defer hb.Stop()
	// simulate service down
	time.AfterFunc(10*time.Second, hb.Stop)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("error accepting:", err)
			return
		}
		log.Println("accept conn from:", conn.RemoteAddr())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleRequest(conn)
		}()
	}
}

func handleRequest(conn net.Conn) {
	// do something
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		log.Println("read:", string(buf[:n]))
	}
}
