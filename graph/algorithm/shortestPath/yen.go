package shortestPath

import (
	"fmt"
	"goContainer/graph"
	"goContainer/heap/bheap"
	"math"
)

// YenKSP Yen's k shortest paths algorithm
//	https://en.wikipedia.org/wiki/Yen%27s_algorithm
//	it computes single-source K-shortest loop-less paths for a graph with non-negative edge cost
func YenKSP(g graph.Graph, source, target graph.ID, k int) (dists []float64, paths [][]graph.ID, err error) {
	dists = make([]float64, k)
	for i := 0; i < k; i++ {
		dists[i] = math.Inf(1)
	}
	paths = make([][]graph.ID, k)
	gCopy := g.Copy()

	// 1. Determine the shortest path from the source to the target
	path0, distMap0, err := AStar(gCopy, source, target)
	if err != nil {
		return nil, nil, err
	}
	paths[0] = path0
	dists[0] = distMap0[target]

	// priority queue to store the potential kth shortest path
	QpathK := newPathQ()
	QpathK.push(paths[0], dists[0])
	for i := 1; i <= k; i++ {
		// the spur node ranges from the first node to the next to last node in the previous k-shortest path
		for j := 0; j < len(paths[i-1])-1; j++ {
			// spur node is retrieved from the previous k-shortest path, i - 1
			spurNode := paths[i-1][j]
			// the sequence of nodes from the source to the spur node of the previous k-shortest path
			rootPath := paths[i-1][:j+1]

			fmt.Printf("spurNode: %s, rootPath: %+v\n", spurNode, rootPath)

			for m := 0; m < k; m++ {
				if isSharedRootPath(rootPath, paths[m]) {
					// remove the links that are part of the previous shortest paths which share the same root path
					fmt.Println("deleted edge: ", paths[m][j], paths[m][j+1])
					err = gCopy.DeleteEdge(paths[m][j], paths[m][j+1])
					if err != nil {
						return nil, nil, err
					}
				}
			}

			// remove rootPathNode from Graph without spurNode and source
			// paths[i - 1][:j]
			for _, node := range rootPath[:len(rootPath)-1] {
				if node == source {
					continue
				}
				gCopy.DeleteNode(node)
				fmt.Println("deleted node: ", node)
			}

			// calculate the spur path from the spur node to the target
			spurPath, distMapSpur, err := AStar(gCopy, spurNode, target)
			fmt.Println("spurPath: ", spurPath)
			if err != nil {
				return nil, nil, err
			}

			if distMapSpur[target] != math.Inf(1) {
				// entire path is made up of the root path and spur path
				totalPath := mergePath(rootPath[:len(rootPath)-1], spurPath)
				distSpur := getPathAccumulativeDistance(totalPath, gCopy)
				fmt.Println("totalPath: ", totalPath, distSpur)
				// add the potential k-shortest path to the heap
				QpathK.push(totalPath, distSpur)
			}

			// add back the edges and nodes that were removed from the graph
			gCopy = g.Copy()
			fmt.Println()
		}

		if QpathK.empty() {
			// this handles the case of there being no spur paths, or no spur paths left.
			// this could happen if the spur paths have already been exhausted (added to `paths`),
			// or there are no spur paths at all - such as when both the source and target vertices
			// lie along a "dead end"
			break
		}

		// add the lowest cost path becomes the k-shortest path
		path, dist := QpathK.pop()
		paths[i-1] = path
		dists[i-1] = dist
	}
	return
}

func getPathAccumulativeDistance(path []graph.ID, g graph.Graph) float64 {
	if len(path) == 0 {
		return math.Inf(1)
	}

	res := 0.

	for i := 0; i < len(path)-1; i++ {
		if _, found := g.GetNode(path[i]); !found {
			fmt.Println("not found: ", path[i])
			return math.Inf(1)
		}
		if edge, err := g.GetEdge(path[i], path[i+1]); err == nil {
			res += edge.Weight()
		} else {
			return math.Inf(1)
		}
	}
	return res
}

func isSharedRootPath(rootPath, path []graph.ID) bool {
	if len(path) < len(rootPath) {
		return false
	}

	return pathEqual(path[:len(rootPath)], rootPath)
}

func pathEqual(pathA, pathB []graph.ID) bool {
	if len(pathA) != len(pathB) {
		return false
	}
	for idx, node := range pathA {
		if pathB[idx] != node {
			return false
		}
	}
	return true
}

func mergePath(pathA, pathB []graph.ID) []graph.ID {
	var nPath []graph.ID
	nPath = append(nPath, pathA...)
	nPath = append(nPath, pathB...)
	return nPath
}

type pathQueue struct {
	items *bheap.MinHeap
}

func pathItemComparator(a, b interface{}) int {
	itemA, itemB := a.(*pathItem), b.(*pathItem)
	A, B := itemA.cost, itemB.cost
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

func (pq *pathQueue) push(path []graph.ID, cost float64) {
	pq.items.Push(newPathItem(path, cost))
}

func (pq *pathQueue) pop() (path []graph.ID, cost float64) {
	it := pq.items.Pop().(*pathItem)
	return it.path, it.cost
}

func (pq *pathQueue) String() string {
	res := ""
	for _, v := range pq.items.Values() {
		it := v.(*pathItem)
		res += fmt.Sprintf("{[%+v], %f}", it.path, it.cost)
	}
	return res
}

func newPathQ() *pathQueue {
	minH := bheap.MinHeap{
		Comparator: pathItemComparator,
	}
	minH.Init()
	return &pathQueue{items: &minH}
}

func (pq *pathQueue) empty() bool {
	return pq.items.Empty()
}

// item data struct used in queue
type pathItem struct {
	path []graph.ID
	// priority
	cost float64
}

func newPathItem(path []graph.ID, cost float64) *pathItem {
	return &pathItem{
		path: path,
		cost: cost,
	}
}
