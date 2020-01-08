package scc

import (
	"goContainer/basic/datastructure/graph"
	"testing"
)

func TestTarjan(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph_scc")
	if err != nil {
		panic(err)
	}

	scc := Tarjan(g)

	if len(scc) != 4 {
		t.Fatalf("expected scc length: 4, Tarjan scc length: %d\n", len(scc))
	}

	expectedSCC := map[int]map[string]struct{}{
		0: {
			"E": {},
			"J": {},
		},

		1: {
			"I": {},
		},

		2: {
			"C": {},
			"D": {},
			"H": {},
		},

		3: {
			"A": {},
			"B": {},
			"F": {},
			"G": {},
		},
	}

	for i, c := range scc {
		for _, cc := range c {
			if _, found := expectedSCC[i][cc.String()]; !found {
				t.Fatalf("%s not found in expected result\n", cc.String())
			}
		}
	}
}
