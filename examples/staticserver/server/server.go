package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Serve ...
// local dev purpose only, do NOT use it in production env
func Serve(dirPath string, addr string) {
	s := http.FileServer(http.Dir(dirPath))
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.ServeHTTP(w, r)
	}
	http.HandleFunc("/", handleFunc)
	fp, _ := filepath.Abs(dirPath)
	log.Println("static server listening at:", addr)
	log.Println("serving dir is:", fp)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
