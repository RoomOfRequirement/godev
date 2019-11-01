package shortestPath

import (
	"goContainer/graph"
)

// GreedyBestFirstSearch (heuristic search)
//	https://en.wikipedia.org/wiki/Best-first_search
//	use a priority queue (minimum heap)
//	compared with Dijkstra algorithm, no priority update process
func GreedyBestFirstSearch(g graph.Graph, source, target graph.ID) ([]graph.ID, map[graph.ID]float64, error) {
	// Q (dist queue): set of all vertices
	Q := newPQFib()
	// prev vertex map
	prev := make(map[graph.ID]graph.ID)
	// distance map
	dist := make(map[graph.ID]float64)
	// initialize
	prev[source] = nil
	Q.push(&item{
		id:   source,
		dist: 0.0,
	})
	dist[source] = 0.

	for !Q.empty() {
		// frontier
		u := Q.pop()

		if u.id == target {
			break
		}

		// next
		tmap, err := g.GetTargets(u.id)
		if err != nil {
			return nil, nil, err
		}

		for v := range tmap {
			// update path
			if _, found := prev[v]; !found {
				e, err := g.GetEdge(u.id, v)
				if err != nil {
					return nil, nil, err
				}
				// nd := heuristic(target, v)
				// since my sample graph only has vertices and edges, no other information like vertex positions (coordinates) etc.
				// the heuristic function can NOT be produced...
				// so here i just use a distance accumulation function to lead a search direction,
				// which means it will search along the smallest distance direction
				// it may be wrong (trap by a local minimum) to find the shortest path (global minimum)
				nd := u.dist + e.Weight()
				Q.push(&item{
					id:   v,
					dist: nd,
				})
				dist[v] = nd
				prev[v] = u.id
			}
		}
	}

	// path list
	var path []graph.ID

	// from target to source
	u := target

	// while prev[u] is defined:
	for {
		if _, existed := prev[u]; !existed {
			break
		}
		// insert u to the beginning of path
		path = append(path, u)
		u = prev[u]
	}

	// reverse path to get the path from source to target
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, dist, nil
}
