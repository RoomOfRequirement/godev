package topologicalsort

import (
	"fmt"
	"goContainer/graph"
	"testing"
)

func TestKahn(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph_topo")
	if err != nil {
		panic(err)
	}
	fmt.Println(g)

	sortedIDs, err := Kahn(g)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(sortedIDs)
	// ABC  DE  FGH
	expected0 := map[graph.ID]struct{}{
		graph.StringID("A"): {},
		graph.StringID("B"): {},
		graph.StringID("C"): {},
	}
	expected1 := map[graph.ID]struct{}{
		graph.StringID("D"): {},
		graph.StringID("E"): {},
	}
	expected2 := map[graph.ID]struct{}{
		graph.StringID("F"): {},
		graph.StringID("G"): {},
		graph.StringID("H"): {},
	}

	for _, id := range sortedIDs[:3] {
		if _, found := expected0[id]; !found {
			t.Fatalf("expected %s not found", id)
		}
	}
	for _, id := range sortedIDs[3:5] {
		if _, found := expected1[id]; !found {
			t.Fatalf("expected %s not found", id)
		}
	}
	for _, id := range sortedIDs[5:] {
		if _, found := expected2[id]; !found {
			t.Fatalf("expected %s not found", id)
		}
	}
}

func TestDFSTopo(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../../test.json", "graph_topo")
	if err != nil {
		panic(err)
	}
	fmt.Println(g)

	sortedIDs, err := DFSTopo(g)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(sortedIDs)
	// AB -> D, BC -> E, ABD -> F, ABCDE -> G, ABDC -> H
	idxMap := make(map[graph.ID]int)
	for i, v := range sortedIDs {
		idxMap[v] = i
	}
	if idxMap[graph.StringID("D")] < idxMap[graph.StringID("A")] || idxMap[graph.StringID("D")] < idxMap[graph.StringID("B")] {
		t.Fatalf("D should behind A and B")
	}
	if idxMap[graph.StringID("E")] < idxMap[graph.StringID("B")] || idxMap[graph.StringID("E")] < idxMap[graph.StringID("C")] {
		t.Fatalf("E should behind B and C")
	}
	if idxMap[graph.StringID("F")] < idxMap[graph.StringID("D")] {
		t.Fatalf("F should behind D")
	}
	if idxMap[graph.StringID("G")] < idxMap[graph.StringID("D")] || idxMap[graph.StringID("G")] < idxMap[graph.StringID("E")] {
		t.Fatalf("G should behind D and E")
	}
	if idxMap[graph.StringID("H")] < idxMap[graph.StringID("C")] || idxMap[graph.StringID("H")] < idxMap[graph.StringID("D")] {
		t.Fatalf("H should behind C and D")
	}
}
