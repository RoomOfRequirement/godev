package service

import (
	"godev/examples/servicediscovery/client"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

// HeartBeat ...
type HeartBeat struct {
	Masters         []string // url
	Name            string   // path = /name
	ServicePort     int
	RequestInterval int64

	StopChan chan struct{}
}

// NewHeartBeat ...
func NewHeartBeat(name string, servicePort int, master string) *HeartBeat {
	return &HeartBeat{
		Masters:         []string{master},
		Name:            name,
		ServicePort:     servicePort,
		RequestInterval: 10,
		StopChan:        make(chan struct{}, 1),
	}
}

// Start ...
func (hb *HeartBeat) Start() {
	go func() {
		for {
			for _, master := range hb.Masters {
				values := make(url.Values)
				values.Add("name", hb.Name)
				values.Add("port", strconv.Itoa(hb.ServicePort))
				_, err := client.Post("http://"+master+"/join", values)
				if err != nil {
					log.Println(err)
				}
			}
			select {
			case <-hb.StopChan:
				log.Println(hb.Name, "stop heart beat")
				return
			default:
				time.Sleep(time.Duration(rand.Int63n(hb.RequestInterval/2)+hb.RequestInterval/2) * time.Second)
			}
		}
	}()
}

// Stop ...
func (hb *HeartBeat) Stop() {
	hb.StopChan <- struct{}{}
}
