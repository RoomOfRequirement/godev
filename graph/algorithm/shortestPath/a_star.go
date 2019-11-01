package shortestPath

import (
	"goContainer/graph"
	"goContainer/heap/bheap"
)

// AStar algorithm
//	https://en.wikipedia.org/wiki/A*_search_algorithm
//	https://www.redblobgames.com/pathfinding/a-star/introduction.html
func AStar(g graph.Graph, source, target graph.ID) ([]graph.ID, map[graph.ID]float64, error) {
	// Q (dist queue): set of all vertices
	//	based on heap: https://github.com/Harold2017/golina/tree/master/container/heap/bheap
	//	use this heap because fibonacci heap `Update` method (based on `decreaseKey` / `increaseKey`) takes large overhead
	Q := newPQ()
	// prev vertex map
	prev := make(map[graph.ID]graph.ID)
	// distance map
	dist := make(map[graph.ID]float64)
	// initialize
	prev[source] = nil
	Q.push(source, 0.0)
	dist[source] = 0.

	for !Q.empty() {
		// frontier
		uid, udist := Q.pop()

		if uid == target {
			break
		}

		// next
		tmap, err := g.GetTargets(uid)
		if err != nil {
			return nil, nil, err
		}

		for v := range tmap {
			e, err := g.GetEdge(uid, v)
			if err != nil {
				return nil, nil, err
			}
			nd := udist + e.Weight()

			_, found := prev[v]
			if !found || nd < dist[v] {
				dist[v] = nd
				Q.push(v, nd+heuristic(target, v))
				prev[v] = uid
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

// since my sample graph only has vertices and edges, no other information like vertex positions (coordinates) etc.
// the heuristic function can NOT be produced... just let h(x) = 0
// then it will become dijkstra
func heuristic(target, to graph.ID) float64 {
	return 0
}

type pQueue struct {
	items *bheap.MinHeap
}

func itemComparator(a, b interface{}) int {
	itemA, itemB := a.(*item), b.(*item)
	A, B := itemA.dist, itemB.dist
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

func (pq *pQueue) push(id graph.ID, dist float64) {
	pq.items.Push(newItem(id, dist))
}

func (pq *pQueue) pop() (id graph.ID, dist float64) {
	it := pq.items.Pop().(*item)
	return it.id, it.dist
}

func newPQ() *pQueue {
	minH := bheap.MinHeap{
		Comparator: itemComparator,
	}
	minH.Init()
	return &pQueue{items: &minH}
}

func (pq *pQueue) empty() bool {
	return pq.items.Empty()
}
