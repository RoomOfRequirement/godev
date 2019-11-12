package shortestpath

import (
	"fmt"
	"goContainer/graph"
	"testing"
)

func TestYenKSP(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph_yen")
	if err != nil {
		panic(err)
	}
	fmt.Println("original graph:")
	fmt.Println(g)
	distance, path, err := YenKSP(g, graph.StringID("C"), graph.StringID("H"), 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
	fmt.Println(distance)
	// C E F H  5
	// C E G H  7
	// C D F H  8
	expectedDistPath := map[float64][]graph.ID{
		5: {graph.StringID("C"), graph.StringID("E"), graph.StringID("F"), graph.StringID("H")},
		7: {graph.StringID("C"), graph.StringID("E"), graph.StringID("G"), graph.StringID("H")},
		8: {graph.StringID("C"), graph.StringID("D"), graph.StringID("F"), graph.StringID("H")},
	}
	for i, dist := range distance {
		if !pathEqual(expectedDistPath[dist], path[i]) {
			t.Fail()
		}
	}
}
