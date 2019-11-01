package graph

// ID: generic interface for any id
type ID interface {
	String() string
}

// Node: generic interface for node or vertex in graph
type Node interface {
	ID() ID
	String() string
}

// Edge: edge includes two nodes (source, target) and its corresponding weight from source points to target
type Edge interface {
	Source() Node
	Target() Node
	Weight() float64
	String() string
}

// EdgeSlice: for sorting
type EdgeSlice []Edge

func (es EdgeSlice) Len() int {
	return len(es)
}

func (es EdgeSlice) Less(i, j int) bool {
	return es[i].Weight() < es[j].Weight()
}

func (es EdgeSlice) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

// Graph: generic interface for using hash-map to store nodes to avoid duplicates
type Graph interface {
	NodeNum() int
	EdgeNum() int

	AddNode(node Node) bool
	DeleteNode(id ID) bool
	// error if id does NOT exist
	ReplaceNode(id ID, newNode Node) error
	GetNode(id ID) (node Node, existed bool)

	// error if source or target node does NOT exist or no edge between
	AddEdge(idSource, idTarget ID, weight float64) error
	DeleteEdge(idSource, idTarget ID) error
	ReplaceEdge(idSource, idTarget ID, weight float64) error
	GetEdge(idSource, idTarget ID) (Edge, error)

	// use hash-map to store nodes, then no duplicate is allowed
	GetNodes() map[ID]Node
	// get all nodes which point to current node
	GetSources(id ID) (map[ID]Node, error)
	// get all nodes which current node points to
	GetTargets(id ID) (map[ID]Node, error)
	// get all edges of current node
	GetEdges(id ID) (inEdges, outEdges EdgeSlice, err error)
	// get all edges which fan in current node
	GetInEdges(id ID) (EdgeSlice, error)
	// get all edges which current node fans out
	GetOutEdges(id ID) (EdgeSlice, error)

	String() string

	// following to meet container interface
	Size() int // here define size = node number
	Empty() bool
	Clear()
	Values() []interface{} // here define value = node

	Copy() Graph // return a deep copy of graph
}
