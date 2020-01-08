package topologicalsort

import (
	"fmt"
	"goContainer/basic/datastructure/graph"
	"goContainer/basic/datastructure/queue/deque"
)

// Kahn algorithm
//	https://en.wikipedia.org/wiki/Topological_sorting
func Kahn(g graph.Graph) (sortedIDs []graph.ID, err error) {
	// 1. compute in-edges for every vertex and initialize visited vertices number as 0
	nodeWithInDegree := make(map[graph.ID]int, g.NodeNum())
	for id := range g.GetNodes() {
		sources, err := g.GetSources(id)
		if err != nil {
			return nil, err
		}
		nodeWithInDegree[id] = len(sources)
	}
	visited := 0

	// 2. enqueue all vertices with degree 0
	Q := deque.NewDeque(g.NodeNum())
	for id, inDegree := range nodeWithInDegree {
		if inDegree == 0 {
			Q.PushBack(id)
		}
	}

	// 3. dequeue a vertex until Q is empty
	//	a. increment visited by 1
	//	b. decrease all vertex's neighbors' degree by 1
	//	c. if in-degree of a neighbor is reduced to 0, then enqueue this neighbor
	for !Q.Empty() {
		v, err := Q.PopFront()
		if err != nil {
			return nil, err
		}
		id := v.(graph.ID)
		sortedIDs = append(sortedIDs, id)

		visited++

		tmap, err := g.GetTargets(id)
		if err != nil && err.Error() != graph.NodeNotExistError(id).Error() {
			return nil, err
		}

		for t := range tmap {
			nodeWithInDegree[t]--
			if nodeWithInDegree[t] == 0 {
				Q.PushBack(t)
			}
		}
	}

	if visited != g.NodeNum() {
		return nil, fmt.Errorf("graph is not a DAG, can NOT do topological sort on it")
	}
	return
}

// DFSTopo depth-first search
//	it sounds like tri-color marking algorithm (golang GC algorithm) but different
//	https://en.wikipedia.org/wiki/Tracing_garbage_collection#Tri-color_marking
//	https://gist.github.com/Harold2017/7529971396e09992f879b22663726e07
func DFSTopo(g graph.Graph) (sortedIDs []graph.ID, err error) {
	// unmarked: 0, temporary mark: 1, permanent mark 2
	mark := make(map[graph.ID]int, g.NodeNum())
	// unmark all vertices
	for v := range g.GetNodes() {
		mark[v] = 0
	}
	// recursively
	for v := range g.GetNodes() {
		if mark[v] == 0 {
			err := visit(g, v, &sortedIDs, &mark)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

func visit(g graph.Graph, id graph.ID, sortedIDs *[]graph.ID, mark *map[graph.ID]int) error {
	if (*mark)[id] == 2 {
		return nil
	}
	if (*mark)[id] == 1 {
		return fmt.Errorf("graph is not a DAG, can NOT do topological sort on it")
	}
	(*mark)[id] = 1
	tmap, err := g.GetTargets(id)
	if err != nil && err.Error() != graph.NodeNotExistError(id).Error() {
		return err
	}
	for t := range tmap {
		err := visit(g, t, sortedIDs, mark)
		if err != nil {
			return err
		}
	}
	(*mark)[id] = 2
	*sortedIDs = append([]graph.ID{id}, *sortedIDs...)
	return nil
}

// TODO: how to implement parallel algorithm in wiki?
