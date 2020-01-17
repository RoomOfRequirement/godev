package traversal

import (
	"fmt"
	"godev/basic/datastructure/graph"
	"testing"
)

func TestDFS(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph")
	if err != nil {
		panic(err)
	}

	nodes, err := DFS(g, graph.StringID("A"))

	if err != nil || len(nodes) != 8 {
		fmt.Println(err, len(nodes))
		t.Fail()
	}

	// fmt.Println(nodes)
}
