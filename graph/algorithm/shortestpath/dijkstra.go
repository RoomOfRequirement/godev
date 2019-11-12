package shortestpath

import (
	"goContainer/graph"
	"goContainer/heap"
	"goContainer/heap/fibonacci"
	"goContainer/queue/deque"
	"math"
)

// Dijkstra algorithm
//	https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
//	use priority queue (minimum heap) (fibonacci heap: https://github.com/Harold2017/goContainer/tree/master/queue/prque/pqfibo)
func Dijkstra(g graph.Graph, source, target graph.ID) ([]graph.ID, map[graph.ID]float64, error) {
	// Q (dist queue): set of all vertices
	Q := newPQFib()
	// map to store all distances
	dist := make(map[graph.ID]*item)
	dist[source] = &item{
		id:   source,
		dist: 0.0,
	}

	// initialize
	// loop all vertices
	for id := range g.GetNodes() {
		if id != source {
			// set initial d to inf
			dist[id] = newItem(id, math.MaxFloat64)
		}

		// push vertex with dist into Q
		it := dist[id]
		Q.push(it)
	}

	// prev vertex map
	prev := make(map[graph.ID]graph.ID)

	for !Q.empty() {
		// extract item with min dist
		u := Q.pop()

		if u.id == target {
			break
		}

		// loop all children vertices of u
		tmap, err := g.GetTargets(u.id)
		if err != nil {
			return nil, nil, err
		}

		for v := range tmap {
			// update distance
			e, err := g.GetEdge(u.id, v)
			if err != nil {
				return nil, nil, err
			}
			nd := dist[u.id].dist + e.Weight()
			// compare with v
			if dist[v].dist > nd {
				// update distance (this will update Q as well)
				it := dist[v]
				nit := &item{
					id:   it.id,
					dist: nd,
				}
				dist[v] = nit
				// update prev vertex map
				prev[v] = u.id
				// update Q
				Q.update(it, nit)
			}
		}
	}

	// path
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

	// dist to pDist
	pDist := make(map[graph.ID]float64, len(dist))
	for k, it := range dist {
		pDist[k] = it.dist
	}

	return path, pDist, nil
}

type priorityQueue struct {
	dists *fibonacci.Heap
}

func newPQFib() *priorityQueue {
	return &priorityQueue{dists: fibonacci.NewHeap()}
}

func (pq *priorityQueue) push(it *item) {
	pq.dists.Insert(it)
}

func (pq *priorityQueue) pop() *item {
	return pq.dists.DeleteMin().(*item)
}

func (pq *priorityQueue) update(it, nit *item) {
	pq.dists.Update(it, nit)
}

func (pq *priorityQueue) empty() bool {
	return pq.dists.Empty()
}

// item data struct used in queue
type item struct {
	id graph.ID
	// priority
	dist float64
}

// Compare to meet heap.Item interface
func (it *item) Compare(ait heap.Item) int {
	if it.dist > ait.(*item).dist {
		return 1
	} else if it.dist == ait.(*item).dist {
		return 0
	} else {
		return -1
	}
}

func newItem(id graph.ID, dist float64) *item {
	return &item{
		id:   id,
		dist: dist,
	}
}
