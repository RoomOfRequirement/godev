package shortestpath

import (
	"fmt"
	"goContainer/basic/datastructure/graph"
	"testing"
)

func TestGreedyBestFirstSearch(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph")
	if err != nil {
		panic(err)
	}
	path, distance, err := GreedyBestFirstSearch(g, graph.StringID("A"), graph.StringID("E"))
	if err != nil {
		panic(err)
	}
	var ts []string
	for _, v := range path {
		ts = append(ts, fmt.Sprintf("%s(%.2f)", v, distance[v]))
	}
	fmt.Println(ts)

	if len(path) != 3 || path[1] != graph.StringID("B") {
		t.Fail()
	}
	// notice the result is  different from Dijkstra algorithm
}
