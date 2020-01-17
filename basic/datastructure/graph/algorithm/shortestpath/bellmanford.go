package shortestpath

import (
	"fmt"
	"godev/basic/datastructure/graph"
	"godev/basic/datastructure/queue/deque"
	"math"
)

// BellmanFord algorithm
//	https://en.wikipedia.org/wiki/Bellman%E2%80%93Ford_algorithm
//	BellmanFord algorithm can work with negative weight edges
//	Complexity is O(|V| * |E|)
func BellmanFord(g graph.Graph, source, target graph.ID) ([]graph.ID, map[graph.ID]float64, error) {
	// map to store all distances
	dist := make(map[graph.ID]float64)
	dist[source] = 0.

	// Step 1: initialize graph
	// for every vertex in graph
	for v := range g.GetNodes() {
		// if vertex not source
		if v != source {
			// distance = Inf
			dist[v] = math.MaxFloat64
		}
	}

	prev := make(map[graph.ID]graph.ID)

	// Step 2: relax edges repeatedly
	// for 1 to |V| - 1
	for i := 1; i < g.NodeNum(); i++ {
		// for every Edge(u, v)
		for u := range g.GetNodes() {
			tmap, err := g.GetTargets(u)
			if err != nil {
				return nil, nil, err
			}
			for v := range tmap {
				// Edge(u, v)
				e, err := g.GetEdge(u, v)
				if err != nil {
					return nil, nil, err
				}

				nd := dist[u] + e.Weight()
				if dist[v] > nd {
					dist[v] = nd
					// refresh prev
					prev[v] = u
				}
			}

			// bi-directional graph (A <-> B), no need for single directional
			// check Edge(v, u)
			smap, err := g.GetSources(u)
			if err != nil {
				return nil, nil, err
			}

			for v := range smap {
				// Edge(v, u)
				e, err := g.GetEdge(v, u)
				if err != nil {
					return nil, nil, err
				}

				nd := dist[v] + e.Weight()
				if dist[u] > nd {
					dist[u] = nd
					// refresh prev
					prev[u] = v
				}
			}
		}
	}

	// Step 3: check for negative-weight cycles
	// for every Edge(u, v)
	for u := range g.GetNodes() {
		tmap, err := g.GetTargets(u)
		if err != nil {
			return nil, nil, err
		}

		for v := range tmap {
			// Edge(u, v)
			e, err := g.GetEdge(u, v)
			if err != nil {
				return nil, nil, err
			}

			nd := dist[u] + e.Weight()
			if dist[v] > nd {
				// negative cycle
				return nil, nil, fmt.Errorf("there exists negative cycle: %v", g)
			}
		}

		// bi-directional graph (A <-> B), no need for single directional
		// check Edge(v, u)
		smap, err := g.GetSources(u)
		if err != nil {
			return nil, nil, err
		}

		for v := range smap {
			// Edge(v, u)
			e, err := g.GetEdge(v, u)
			if err != nil {
				return nil, nil, err
			}

			nd := dist[v] + e.Weight()
			if dist[u] > nd {
				// negative cycle
				return nil, nil, fmt.Errorf("there exists negative cycle: %v", g)
			}
		}
	}

	// path
	//	for small graph with a, actually no need to use deque,
	//	just use slice copy to realize push_front operation and then no need to do transfer from deque to slice...
	//	anyway, here just for demo the algorithm
	pathQ := deque.NewDeque(0)

	u := target

	// while prev[u] is defined:
	for {
		if _, existed := prev[u]; !existed {
			break
		}
		// insert u to the beginning of path
		pathQ.PushFront(u)
		u = prev[u]
	}

	// add source
	pathQ.PushFront(source)

	// pathQ to []ID
	path := make([]graph.ID, 0, pathQ.Size())
	for !pathQ.Empty() {
		e, err := pathQ.PopFront()
		if err != nil {
			return nil, nil, err
		}
		path = append(path, e.(graph.ID))
	}

	return path, dist, nil
}
