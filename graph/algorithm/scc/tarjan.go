package scc

import (
	"goContainer/graph"
)

// Tarjan algorithm
//	https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
//	DFS based
func Tarjan(g graph.Graph) [][]graph.ID {
	data := newData()

	for id := range g.GetNodes() {
		if _, found := data.index[id]; !found {
			strongConnect(g, id, data)
		}
	}

	return data.result
}

func strongConnect(g graph.Graph, id graph.ID, data *data) {
	// set the depth index for v to the smallest unused index
	data.index[id] = data.globalIndex
	data.lowLink[id] = data.globalIndex
	data.globalIndex++
	// push in stack
	data.stack.push(id)

	// consider successors of v
	tmap, err := g.GetTargets(id)
	if err != nil && err.Error() != graph.NodeNotExistError(id).Error() {
		panic(err)
	}

	for w := range tmap {
		// successor w has not yet been visited; recurse on it
		if _, found := data.index[w]; !found {
			strongConnect(g, w, data)
			data.lowLink[id] = min(data.lowLink[id], data.lowLink[w])
		} else if data.stack.in(w) {
			// successor w is in stack and hence in the current strong connect component (SCC)
			data.lowLink[id] = min(data.lowLink[id], data.index[w])
		}
	}

	// if v is a root node, pop the stack and generate an SCC
	if data.lowLink[id] == data.index[id] {
		// start a new strongly connected component
		var scc []graph.ID
		for {
			u := data.stack.pop()
			scc = append(scc, u)
			if u == id {
				data.result = append(data.result, scc)
				break
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type data struct {
	globalIndex    int
	index, lowLink map[graph.ID]int
	stack          *stack
	result         [][]graph.ID
}

func newData() *data {
	return &data{
		globalIndex: 0,
		index:       make(map[graph.ID]int),
		lowLink:     make(map[graph.ID]int),
		stack:       newStack(),
		result:      [][]graph.ID{},
	}
}

type stack struct {
	// stack
	s []graph.ID
	// map for quick lookup whether id is on stack
	m map[graph.ID]struct{}
}

func newStack() *stack {
	return &stack{
		s: []graph.ID{},
		m: make(map[graph.ID]struct{}),
	}
}

func (st *stack) push(id graph.ID) {
	st.s = append(st.s, id)
	st.m[id] = struct{}{}
}

func (st *stack) pop() graph.ID {
	v := st.s[len(st.s)-1]
	st.s = st.s[:len(st.s)-1]
	delete(st.m, v)
	return v
}

func (st *stack) in(id graph.ID) bool {
	_, found := st.m[id]
	return found
}
