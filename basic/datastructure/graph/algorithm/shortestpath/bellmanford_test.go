package shortestpath

import (
	"fmt"
	"goContainer/basic/datastructure/graph"
	"testing"
)

func TestBellmanFord(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph_n")
	if err != nil {
		panic(err)
	}
	path, distance, err := BellmanFord(g, graph.StringID("E"), graph.StringID("D"))
	if err != nil {
		panic(err)
	}
	var ts []string
	for _, v := range path {
		ts = append(ts, fmt.Sprintf("%s(%.2f)", v, distance[v]))
	}
	fmt.Println(ts)

	if len(path) != 5 {
		t.Fail()
	}

	expectedPath := []string{"E", "A", "C", "B", "D"}
	for i, v := range ts {
		if string(v[0]) != expectedPath[i] {
			fmt.Println(i, v)
			t.Fail()
		}
	}

	// negative cycle E(0.00) -> B(6.00) -> D(2.00) -> E(-1.00)
	g, err = graph.NewGraphFromJSON("../../test.json", "graph_nc")
	if err != nil {
		panic(err)
	}
	path, distance, err = BellmanFord(g, graph.StringID("E"), graph.StringID("D"))
	if path != nil || distance != nil || err == nil {
		t.Fail()
	}
}
