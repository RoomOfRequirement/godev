package main

import (
	"flag"
	"fmt"
	"godev/examples/staticserver/server"
	"os"
)

var (
	help       bool
	version    bool
	dirPath    string
	serverAddr string
)

func init() {
	flag.BoolVar(&help, "h", false, "help info")
	flag.BoolVar(&version, "v", false, "version info")
	flag.StringVar(&dirPath, "p", ".", "static dir path")
	flag.StringVar(&serverAddr, "l", ":8888", "listener addr")
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `Simple Static Server to serve your static files
Notice: Only use in local dev env, NOT in production env!
Version: 0.0.1
Usage: ss [-hvsftld] [-h help] [-v version] [-p static dir path] [-l server listening address]
Options
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	} else if version {
		fmt.Println("version: 0.0.1")
	} else {
		server.Serve(dirPath, serverAddr)
	}
}
