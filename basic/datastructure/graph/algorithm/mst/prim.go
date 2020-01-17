package mst

import (
	"godev/basic/datastructure/graph"
	"godev/basic/datastructure/heap"
	"godev/basic/datastructure/heap/fibonacci"
	"math"
)

// Prim algorithm
//	https://en.wikipedia.org/wiki/Prim%27s_algorithm
//	greedy algorithm, need a priority queue
//	similar with dijkstra algorithm (https://github.com/Harold2017/golina/blob/master/container/graph/algorithm/shortest_path/dijkstra.go)
//	since i use a map to record edge info, the result slice has NO order (compared with Kruskal, where edges are sorted)
func Prim(g graph.Graph, startVetex graph.ID) (graph.EdgeSlice, error) {
	Q := newPQFib()
	dist := make(map[graph.ID]float64)
	dist[startVetex] = 0.
	items := make(map[graph.ID]*item)

	for id := range g.GetNodes() {
		if id != startVetex {
			dist[id] = math.MaxFloat64
		}
		it := newItem(id, dist[id])
		Q.push(it)
		items[id] = it
	}

	prev := make(map[graph.ID]graph.ID)

	for !Q.empty() {
		// nearest vertex
		u := Q.pop()

		// adjacent vertices
		// out-edges
		tmap, _ := g.GetTargets(u.id)

		for tid := range tmap {
			if it := items[tid]; Q.search(it) {
				edge, _ := g.GetEdge(u.id, tid)
				if edge == nil {
					continue
				}
				w := edge.Weight()
				if dist[tid] > w {
					dist[tid] = w
					prev[tid] = u.id
					nit := newItem(tid, w)
					Q.update(it, nit)
					items[tid] = nit
				}
			}
		}

		// in-edges
		smap, _ := g.GetSources(u.id)

		for sid := range smap {
			if it := items[u.id]; Q.search(it) {
				edge, _ := g.GetEdge(sid, u.id)
				if edge == nil {
					continue
				}
				w := edge.Weight()
				if dist[u.id] > w {
					dist[u.id] = w
					prev[u.id] = sid
					nit := newItem(u.id, w)
					Q.update(it, nit)
					items[u.id] = nit
				}
			}
		}
	}

	edges := graph.EdgeSlice{}

	for k, v := range prev {
		edge, _ := g.GetEdge(v, k)
		edges = append(edges, edge)
	}
	return edges, nil
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

func (pq *priorityQueue) search(it *item) bool {
	return pq.dists.Search(it)
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
