package reverse

import (
	"godev/controlflow/loadbalancer/roundrobin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type servers struct {
	addresses []string
	weights   []int
}

// New ...
func New(addresses []string, weights []int) http.Handler {
	return &servers{addresses: addresses, weights: weights}
}

// ServeHTTP ...
func (ss *servers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr, _ := roundrobin.New(roundrobin.Nginx)
	_ = rr.SetNodes(ss.addresses, ss.weights)
	addr := rr.Next()
	server, _ := url.Parse("http://" + addr)
	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(w, r)
}
