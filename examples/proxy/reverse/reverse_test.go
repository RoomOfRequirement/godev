package reverse

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type helloHandler struct{}

func (*helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

var upgrader = websocket.Upgrader{}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// echo
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func TestNew(t *testing.T) {
	listen := ":10086"
	s1, s2, s3 := "127.0.0.1:9999", "127.0.0.1:8888", "127.0.0.1:6666"
	w1, w2, w3 := 2, 4, 8
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
	// ws
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", ws)
		s := &http.Server{
			Addr:           ":8888",
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

	time.Sleep(100 * time.Millisecond)
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost" + listen,
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer func() {
		_ = c.Close()
	}()
	err = c.WriteMessage(websocket.TextMessage, []byte("Hello World!"))
	assert.NoError(t, err)
	_, message, err := c.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", string(message))
}
