package server

import (
	"godev/examples/simpleweb/bb"
	"log"
	"net/http"
	"strings"
	"time"
)

// RemoteService ...
type RemoteService struct {
	Name         string
	Location     string // url
	LastHearBeat int64  // seconds
}

// Master ...
type Master struct {
	// service heart beat timeout
	Timeout int64
	// path -> services
	ServiceMapping map[string][]RemoteService // one path can have multiple services
}

func (ms *Master) ping(ctx *bb.Context) {
	ctx.String(http.StatusOK, "hello %s", "world")
}

func (ms *Master) list(ctx *bb.Context) {
	path := ctx.Query("path")
	log.Println(path)
	resp := make([]RemoteService, 0)
	services, found := ms.ServiceMapping[path]
	if !found {
		ctx.JSON(http.StatusOK, bb.OBJ{
			"services": resp,
		})
		return
	}
	now := time.Now().Unix()
	for _, service := range services {
		if service.LastHearBeat+ms.Timeout >= now {
			resp = append(resp, service)
		}
	}
	// update ms
	ms.ServiceMapping[path] = resp
	ctx.JSON(http.StatusOK, bb.OBJ{
		"services": resp,
	})
}

func (ms *Master) join(ctx *bb.Context) {
	name := ctx.PostForm("name")
	path := name // for simplification
	port := ctx.PostForm("port")
	host := ctx.Req.Host
	if strings.Contains(host, ":") {
		host = host[0:strings.Index(host, ":")]
	}
	location := host + ":" + port
	services, found := ms.ServiceMapping[path]
	if !found {
		services = make([]RemoteService, 0)
	}
	found = false
	for _, service := range services {
		if service.Location == location {
			service.LastHearBeat = time.Now().Unix()
			found = true
			break
		}
	}
	if !found {
		services = append(services, RemoteService{
			Name:         name,
			Location:     location,
			LastHearBeat: time.Now().Unix(),
		})
	}
	// update ms
	ms.ServiceMapping[path] = services
	ctx.JSON(http.StatusAccepted, bb.OBJ{
		"services": services,
	})
}

// Run ...
func Run(addr string, timeoutLimit int64) {
	ms := &Master{
		Timeout:        timeoutLimit,
		ServiceMapping: make(map[string][]RemoteService),
	}
	b := bb.Default()
	b.GET("/", ms.ping)
	b.GET("/list", ms.list)
	b.POST("/join", ms.join)
	_ = b.Run(addr)
}
