package traversal

import (
	"fmt"
	"godev/basic/datastructure/graph"
)

// BFS breadth first search
//	https://en.wikipedia.org/wiki/Breadth-first_search
//	time complexity: O(|V| + |E|), O(|E|): O(1) - O(|V|^2)
func BFS(g graph.Graph, startVertex graph.ID) ([]graph.Node, error) {
	startNode, found := g.GetNode(startVertex)
	if !found {
		return nil, fmt.Errorf("start vertex is not found in the graph")
	}

	// queue (FIFO)
	// enqueue statNode
	Q := []graph.Node{startNode}
	// result node array
	res := []graph.Node{startNode}

	// record visited state
	visited := make(map[graph.ID]struct{})
	visited[startNode.ID()] = struct{}{}

	for len(Q) != 0 {
		// dequeue
		v := Q[0]
		Q = Q[1:]

		// loop adjacent nodes of v
		//	adjacent nodes including two parts: targets (out), sources (in)

		// targets
		targets, err := g.GetTargets(v.ID()) // map[ID]Node
		if err != nil {
			return nil, err
		}
		for tID, tNode := range targets {
			// if t not visited
			if _, found := visited[tID]; !found {
				// label t as visited
				visited[tID] = struct{}{}
				// enqueue t
				Q = append(Q, tNode)

				res = append(res, tNode)
			}
		}

		// sources
		sources, err := g.GetSources(v.ID()) // map[ID]Node
		if err != nil {
			return nil, err
		}
		for sID, sNode := range sources {
			// if s not visited
			if _, found := visited[sID]; !found {
				// label s as visited
				visited[sID] = struct{}{}
				// enqueue s
				Q = append(Q, sNode)

				res = append(res, sNode)
			}
		}
	}

	return res, nil
}
