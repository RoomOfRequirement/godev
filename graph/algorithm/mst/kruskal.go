package mst

import (
	"goContainer/graph"
	"sort"
)

// Kruskal algorithm
//	https://en.wikipedia.org/wiki/Kruskal%27s_algorithm
func Kruskal(g graph.Graph) (graph.EdgeSlice, error) {
	dss := NewDSSet()
	for v := range g.GetNodes() {
		s := MakeSet()
		s.Element[v] = struct{}{}
		dss.Sets[s] = struct{}{}
	}

	edges := graph.EdgeSlice{}
	// to remove duplicate
	foundEdge := make(map[string]struct{})
	for id := range g.GetNodes() {
		inEdges, outEdges, _ := g.GetEdges(id)
		if inEdges == nil && outEdges == nil {
			return nil, graph.OrphanNodeError(id)
		}
		if inEdges != nil {
			for _, edge := range inEdges {
				if _, found := foundEdge[edge.String()]; !found {
					edges = append(edges, edge)
					foundEdge[edge.String()] = struct{}{}
				}
			}
		}
		if outEdges != nil {
			for _, edge := range outEdges {
				if _, found := foundEdge[edge.String()]; !found {
					edges = append(edges, edge)
					foundEdge[edge.String()] = struct{}{}
				}
			}
		}
	}

	// sort edges in ascending order according to weight
	sort.Sort(edges)

	A := graph.EdgeSlice{}

	for _, edge := range edges {
		findSetU, findSetV := dss.FindSet(edge.Source().ID()).Find(), dss.FindSet(edge.Target().ID()).Find()
		if findSetU != findSetV {
			A = append(A, edge)
			// union
			Union(findSetU, findSetV)
		}
	}

	return A, nil
}
