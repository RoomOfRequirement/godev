package agent

import (
	"encoding/json"
	"godev/examples/servicediscovery/client"
	"godev/examples/servicediscovery/server"
	"log"
)

// NameServiceAgent ...
type NameServiceAgent struct {
	Masters []string // master url
}

// NewNameServiceAgent ...
func NewNameServiceAgent(masters ...string) *NameServiceAgent {
	return &NameServiceAgent{Masters: masters}
}

// GetService ...
func (nsa *NameServiceAgent) GetService(name string) (locations []string) {
	for _, master := range nsa.Masters {
		jsonBytes, err := client.Get("http://" + master + "/list?path=" + name)
		if err != nil {
			log.Printf("failed to get services list from: %s\n%s\n", master, err)
			continue
		}
		var resp map[string][]server.RemoteService
		err = json.Unmarshal(jsonBytes, &resp)
		if err != nil {
			log.Printf("%s/list?name=%s unmarshal error:%v\njson:%s\n", master, name, err, string(jsonBytes))
			continue
		}
		if len(resp["services"]) < 1 {
			log.Println("agent get no services of name: ", name)
			return nil
		}
		for _, r := range resp["services"] {
			locations = append(locations, r.Location)
		}
	}
	return
}
