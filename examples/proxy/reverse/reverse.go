package reverse

import (
	"godev/controlflow/loadbalancer/roundrobin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type servers struct {
	addresses []string
	weights   []int
	rr        roundrobin.RoundRobin
}

// New ...
func New(addresses []string, weights []int) http.Handler {
	rr, _ := roundrobin.New(roundrobin.Nginx)
	_ = rr.SetNodes(addresses, weights)
	return &servers{addresses: addresses, weights: weights, rr: rr}
}

// ServeHTTP ...
func (ss *servers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := ss.rr.Next()
	var serverURL *url.URL
	prefix := strings.Split(addr, "://")[0]
	switch prefix {
	case "http", "https", "ws", "wss":
		serverURL, _ = url.Parse(addr)
	default:
		serverURL, _ = url.Parse("http://" + addr)
	}
	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.ServeHTTP(w, r)
}
