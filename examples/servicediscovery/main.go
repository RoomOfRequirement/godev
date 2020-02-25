package main

import (
	"godev/examples/servicediscovery/agent"
	"godev/examples/servicediscovery/server"
	"godev/examples/servicediscovery/service"
	"log"
	"net"
	"time"
)

func main() {
	// master
	go server.Run(":9999", 10)
	time.Sleep(time.Second)
	// service
	go service.RunService("test", ":8888", "localhost:9999")
	time.Sleep(time.Second)
	// agent
	ag := agent.NewNameServiceAgent("localhost:9999")
	// client
	for {
		locations := ag.GetService("test")
		for _, location := range locations {
			conn, err := net.Dial("tcp", location)
			if err != nil {
				log.Println(err)
			}
			conn.Write([]byte("hello world!"))
			time.Sleep(time.Second)
			conn.Close()
		}
		time.Sleep(time.Second)
	}
}
