package main

import "godev/examples/staticserver/server"

func main() {
	server.Serve(".", ":8888")
}
