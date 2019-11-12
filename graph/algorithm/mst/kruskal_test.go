package mst

import (
	"fmt"
	"goContainer/graph"
	"testing"
)

func TestKruskal(t *testing.T) {
	g1, err := graph.NewGraphFromJSON("../../test.json", "graph_mst1")
	if err != nil {
		panic(err)
	}
	mst1, err := Kruskal(g1)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	AC, _ := g1.GetEdge(graph.StringID("A"), graph.StringID("C"))
	DE, _ := g1.GetEdge(graph.StringID("D"), graph.StringID("E"))
	AB, _ := g1.GetEdge(graph.StringID("A"), graph.StringID("B"))
	BD, _ := g1.GetEdge(graph.StringID("B"), graph.StringID("D"))
	expected1 := graph.EdgeSlice{AC, DE, AB, BD}
	for i := range mst1 {
		if mst1[i].String() != expected1[i].String() {
			fmt.Println(mst1[i], expected1[i])
			t.Fail()
		}
	}

	g2, err := graph.NewGraphFromJSON("../../test.json", "graph_mst2")
	if err != nil {
		panic(err)
	}
	// mst2 has two order
	// ce == ec in test.json
	mst2, err := Kruskal(g2)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	ce, _ := g2.GetEdge(graph.StringID("C"), graph.StringID("E"))
	ad, _ := g2.GetEdge(graph.StringID("A"), graph.StringID("D"))
	df, _ := g2.GetEdge(graph.StringID("D"), graph.StringID("F"))
	be, _ := g2.GetEdge(graph.StringID("B"), graph.StringID("E"))
	ab, _ := g2.GetEdge(graph.StringID("A"), graph.StringID("B"))
	eg, _ := g2.GetEdge(graph.StringID("E"), graph.StringID("G"))

	ec, _ := g2.GetEdge(graph.StringID("E"), graph.StringID("C"))
	da, _ := g2.GetEdge(graph.StringID("D"), graph.StringID("A"))
	fd, _ := g2.GetEdge(graph.StringID("F"), graph.StringID("D"))
	eb, _ := g2.GetEdge(graph.StringID("E"), graph.StringID("B"))
	ba, _ := g2.GetEdge(graph.StringID("B"), graph.StringID("A"))
	ge, _ := g2.GetEdge(graph.StringID("G"), graph.StringID("E"))
	expected2 := map[string]struct{}{
		ce.String(): {},
		ad.String(): {},
		df.String(): {},
		be.String(): {},
		ab.String(): {},
		eg.String(): {},

		ec.String(): {},
		da.String(): {},
		fd.String(): {},
		eb.String(): {},
		ba.String(): {},
		ge.String(): {},
	}
	for i := range mst2 {
		if _, found := expected2[mst2[i].String()]; !found {
			fmt.Println(mst2[i])
			t.Fail()
		}
	}
}
