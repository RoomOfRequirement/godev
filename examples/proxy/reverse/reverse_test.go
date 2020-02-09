package reverse

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

type helloHandler struct{}

func (*helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func TestNew(t *testing.T) {
	listen := ":10086"
	s1, s2, s3 := "127.0.0.1:6666", "127.0.0.1:8888", "127.0.0.1:9999"
	w1, w2, w3 := 5, 1, 1
	ss := New([]string{s1, s2, s3}, []int{w1, w2, w3})
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", &helloHandler{})
		s := &http.Server{
			Addr:           ":6666",
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		err := s.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		err := http.ListenAndServe(listen, ss)
		if err != nil {
			log.Println(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://" + listen)
	assert.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", string(body))
}
