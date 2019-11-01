package traversal

import (
	"fmt"
	"goContainer/graph"
)

// DFS depth first search
//	https://en.wikipedia.org/wiki/Depth-first_search
//	time complexity: O(|V| + |E|), O(|E|): O(1) - O(|V|^2)
func DFS(g graph.Graph, startVertex graph.ID) ([]graph.Node, error) {
	startNode, found := g.GetNode(startVertex)
	if !found {
		return nil, fmt.Errorf("start vertex is not found in the graph")
	}

	// stack (LIFO)
	// push startNode
	S := []graph.Node{startNode}
	// result node array
	var res []graph.Node

	// record visited state
	visited := make(map[graph.ID]struct{})

	for len(S) != 0 {
		// Pop
		v := S[len(S)-1]
		S = S[:len(S)-1]

		//if v is not visited
		if _, found := visited[v.ID()]; !found {
			// label v as visited
			visited[v.ID()] = struct{}{}

			// update result
			res = append(res, v)

			// loop adjacent nodes of v
			//	adjacent nodes including two parts: targets (out), sources (in)

			// targets
			targets, err := g.GetTargets(v.ID())
			if err != nil {
				return nil, err
			}
			for tID, tNode := range targets {
				// if t is not visited
				if _, found := visited[tID]; !found {
					// Push
					S = append(S, tNode)
				}
			}

			// sources
			sources, err := g.GetSources(v.ID())
			if err != nil {
				return nil, err
			}
			for sID, sNode := range sources {
				// if s is not visited
				if _, found := visited[sID]; !found {
					// Push
					S = append(S, sNode)
				}
			}
		}
	}
	return res, nil
}
