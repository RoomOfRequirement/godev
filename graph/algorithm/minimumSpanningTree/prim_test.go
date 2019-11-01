package minimumSpanningTree

import (
	"fmt"
	"goContainer/graph"
	"testing"
)

func TestPrim(t *testing.T) {
	g1, err := graph.NewGraphFromJSON("../../test.json", "graph_mst1")
	if err != nil {
		panic(err)
	}
	mst1, err := Prim(g1, graph.StringID("A"))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	AC, _ := g1.GetEdge(graph.StringID("A"), graph.StringID("C"))
	DE, _ := g1.GetEdge(graph.StringID("D"), graph.StringID("E"))
	AB, _ := g1.GetEdge(graph.StringID("A"), graph.StringID("B"))
	BD, _ := g1.GetEdge(graph.StringID("B"), graph.StringID("D"))

	expected1 := map[string]struct{}{
		AC.String(): {},
		DE.String(): {},
		AB.String(): {},
		BD.String(): {},
	}
	for i := range mst1 {
		if _, found := expected1[mst1[i].String()]; !found {
			fmt.Println(mst1[i])
			t.Fail()
		}
	}

	g2, err := graph.NewGraphFromJSON("../../test.json", "graph_mst2")
	if err != nil {
		panic(err)
	}
	// mst2 has two order
	// ce == ec in test.json
	mst2, err := Prim(g2, graph.StringID("A"))
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
