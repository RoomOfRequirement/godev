package shortestpath

import (
	"fmt"
	"goContainer/basic/datastructure/graph"
	"testing"
)

func TestAStar(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph")
	if err != nil {
		panic(err)
	}
	path, distance, err := AStar(g, graph.StringID("A"), graph.StringID("E"))
	if err != nil {
		panic(err)
	}
	var ts []string
	for _, v := range path {
		ts = append(ts, fmt.Sprintf("%s(%.2f)", v, distance[v]))
	}
	fmt.Println(ts)

	if len(path) != 3 || path[1] != graph.StringID("D") {
		t.Fail()
	}
}
