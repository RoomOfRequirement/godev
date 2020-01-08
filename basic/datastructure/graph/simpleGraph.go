package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// NodeNotExistError returns error if node not exist
func NodeNotExistError(id ID) error {
	return fmt.Errorf("node id %s does not exist in the graph", id)
}

// EdgeNotExistError returns error if edge not exist
func EdgeNotExistError(idSource, idTarget ID) error {
	return fmt.Errorf("edge from %s to %s does not exist in the graph", idSource, idTarget)
}

// OrphanNodeError returns error if node has no in / out edge
func OrphanNodeError(id ID) error {
	return fmt.Errorf("node %s has no edges fan in/out in the graph", id)
}

// InEdgeNotExistError returns error if node has no in edge
func InEdgeNotExistError(id ID) error {
	return fmt.Errorf("node %s has no edges fan in in the graph", id)
}

// OutEdgeNotExistError returns error if node has no out edge
func OutEdgeNotExistError(id ID) error {
	return fmt.Errorf("node %s has no edges fan out in the graph", id)
}

// StringID alias of customized string ID
type StringID string

// String to match ID interface
func (sid StringID) String() string {
	return string(sid)
}

type node struct {
	id string
}

func (n *node) ID() ID {
	return StringID(n.id)
}

func (n *node) String() string {
	return n.id
}

// NewNode creates a simple node from a string id
func NewNode(id string) Node {
	return &node{id}
}

type edge struct {
	source, target Node
	weight         float64
}

func (e edge) Source() Node {
	return e.source
}

func (e edge) Target() Node {
	return e.target
}

func (e edge) Weight() float64 {
	return e.weight
}

func (e edge) String() string {
	return fmt.Sprintf("%s --> %s (weight: %.6f)\n", e.source, e.target, e.weight)
}

// NewEdge creates a simple edge from two nodes and their edge weight
func NewEdge(source, target Node, weight float64) Edge {
	return &edge{
		source: source,
		target: target,
		weight: weight,
	}
}

type graph struct {
	// sync.RWMutex  // avoid race conditions
	nodes map[ID]Node
	edges map[ID]map[ID]float64 // edges between two nodes with weights (e.g. A B 1.0, B A 0.1)
}

func (g *graph) NodeNum() int {
	return len(g.nodes)
}

func (g *graph) EdgeNum() int {
	return len(g.edges)
}

func (g *graph) AddNode(node Node) bool {
	id := node.ID()
	if _, existed := g.GetNode(id); existed {
		return false
	}
	g.nodes[id] = node
	return true
}

func (g *graph) DeleteNode(id ID) bool {
	if _, existed := g.GetNode(id); !existed {
		return false
	}
	// delete node
	delete(g.nodes, id)

	// delete fan in
	sources, _ := g.GetSources(id)
	for _, m := range sources {
		delete(g.edges[m], id)
	}

	// delete fan out
	delete(g.edges, id)
	return true
}

func (g *graph) ReplaceNode(id ID, newNode Node) error {
	if _, existed := g.GetNode(id); !existed {
		return NodeNotExistError(id)
	}
	// newNode's id may not the same with id
	if newNode.ID() != id {
		return fmt.Errorf("new node id should stay the same with replaced node")
	}
	g.nodes[id] = newNode
	return nil
}

func (g *graph) GetNode(id ID) (node Node, existed bool) {
	node, existed = g.nodes[id]
	return node, existed
}

// add edge from source to target
//	if edge existed, just add new weight to its weight
func (g *graph) AddEdge(idSource, idTarget ID, weight float64) error {
	if _, existed := g.GetNode(idSource); !existed {
		return NodeNotExistError(idSource)
	}
	if _, existed := g.GetNode(idTarget); !existed {
		return NodeNotExistError(idSource)
	}

	if _, existed := g.edges[idSource]; existed {
		if _, _existed := g.edges[idSource][idTarget]; _existed {
			// if edge existed, just update weight by adding
			g.edges[idSource][idTarget] += weight
		} else {
			// else, set edge with weight
			g.edges[idSource][idTarget] = weight
		}
	} else {
		g.edges[idSource] = make(map[ID]float64)
		g.edges[idSource][idTarget] = weight
	}
	return nil
}

func (g *graph) DeleteEdge(idSource, idTarget ID) error {
	if _, existed := g.GetNode(idSource); !existed {
		return NodeNotExistError(idSource)
	}
	if _, existed := g.GetNode(idTarget); !existed {
		return NodeNotExistError(idSource)
	}

	if e, existed := g.edges[idSource]; existed {
		if _, existed := e[idTarget]; existed {
			delete(g.edges[idSource], idTarget)
		}
	}
	// EdgeNotExistError(idSource, idTarget)
	return nil
}

// if edge existed, replace edge weight
//	else create edge with weight
func (g *graph) ReplaceEdge(idSource, idTarget ID, weight float64) error {
	if _, existed := g.GetNode(idSource); !existed {
		return NodeNotExistError(idSource)
	}
	if _, existed := g.GetNode(idTarget); !existed {
		return NodeNotExistError(idSource)
	}

	if _, existed := g.edges[idSource]; existed {
		g.edges[idSource][idTarget] = weight
	} else {
		g.edges[idSource] = make(map[ID]float64)
		g.edges[idSource][idTarget] = weight
	}
	return nil
}

func (g *graph) GetEdge(idSource, idTarget ID) (Edge, error) {
	if _, existed := g.GetNode(idSource); !existed {
		return nil, NodeNotExistError(idSource)
	}
	if _, existed := g.GetNode(idTarget); !existed {
		return nil, NodeNotExistError(idSource)
	}
	weight, existed := g.edges[idSource][idTarget]
	if !existed {
		return nil, EdgeNotExistError(idSource, idTarget)
	}
	return NewEdge(g.nodes[idSource], g.nodes[idTarget], weight), nil
}

func (g *graph) GetNodes() map[ID]Node {
	return g.nodes
}

// TODO: this function need to transverse the whole edges map now
//  whether need to add an additional map to store sources,
//  which will improve this lookup but introduce more operations on add/delete
func (g *graph) GetSources(id ID) (map[ID]Node, error) {
	if _, existed := g.GetNode(id); !existed {
		return nil, NodeNotExistError(id)
	}
	sources := map[ID]Node{}
	for i, e := range g.edges {
		if _, existed := e[id]; existed {
			sources[i] = g.nodes[i]
		}
	}
	return sources, nil
}

func (g *graph) GetTargets(id ID) (map[ID]Node, error) {
	nodes, existed := g.edges[id]
	if !existed {
		return nil, NodeNotExistError(id)
	}
	rs := map[ID]Node{}
	for n := range nodes {
		rs[n] = g.nodes[n]
	}
	return rs, nil
}

func (g *graph) GetEdges(id ID) (fanIn, fanOut EdgeSlice, err error) {
	if _, existed := g.nodes[id]; !existed {
		return nil, nil, NodeNotExistError(id)
	}
	fanIn, err = g.GetInEdges(id)
	if err != nil {
		return
	}
	fanOut, err = g.GetOutEdges(id)
	return
}

func (g *graph) GetInEdges(id ID) (EdgeSlice, error) {
	sources, err := g.GetSources(id)
	if err != nil {
		return nil, err
	}
	if sources == nil {
		return nil, InEdgeNotExistError(id)
	}
	es := make(EdgeSlice, 0, len(sources))
	for _, t := range sources {
		e, _ := g.GetEdge(t.ID(), id)
		es = append(es, e)
	}
	return es, nil
}

func (g *graph) GetOutEdges(id ID) (EdgeSlice, error) {
	targets, err := g.GetTargets(id)
	if err != nil {
		return nil, err
	}
	if targets == nil {
		return nil, OutEdgeNotExistError(id)
	}
	es := make(EdgeSlice, 0, len(targets))
	for _, t := range targets {
		e, _ := g.GetEdge(id, t.ID())
		es = append(es, e)
	}
	return es, nil
}

func (g *graph) String() string {
	buf := new(bytes.Buffer)
	for id, source := range g.nodes {
		targets, _ := g.GetTargets(id)
		for _, nodeTarget := range targets {
			weight := g.edges[id][nodeTarget.ID()]
			// ignore error
			_, _ = fmt.Fprintf(buf, "%s --> %s (weight: %.6f)\n", source, nodeTarget, weight)
		}
	}
	return buf.String()
}

func (g *graph) Size() int {
	return g.NodeNum()
}

func (g *graph) Empty() bool {
	return g.Size() == 0
}

func (g *graph) Clear() {
	g.nodes = make(map[ID]Node)
	g.edges = make(map[ID]map[ID]float64)
}

func (g *graph) Values() []interface{} {
	values := make([]interface{}, 0, g.NodeNum())
	for _, n := range g.nodes {
		values = append(values, n.ID())
	}
	return values
}

// NewGraph returns a simple graph which meets the Graph interface
func NewGraph() Graph {
	return &graph{
		// RWMutex: sync.RWMutex{},
		nodes: make(map[ID]Node),
		edges: make(map[ID]map[ID]float64),
	}
}

/*
NewGraphFromJSON creates a simple graph from json file which has the following structure:

{
	"graph": {
		"A": {
			"B": 1,
			"C": 2
		},
		"B": {
			"A": 3
		},
		"C": {
			"B": 5
		}
	}
}
*/
func NewGraphFromJSON(filePath, graphName string) (Graph, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	mapping := map[string]map[string]map[string]float64{}
	decoder := json.NewDecoder(file)
	for {
		if err := decoder.Decode(&mapping); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if _, existed := mapping[graphName]; !existed {
		return nil, fmt.Errorf("%s does not exist", graphName)
	}
	graph := NewGraph()
	for id, m := range mapping[graphName] {
		// check duplicates
		if node, existed := graph.GetNode(StringID(id)); !existed {
			node = NewNode(id)
			graph.AddNode(node)
		}
		for id2, weight := range m {
			if node2, existed := graph.GetNode(StringID(id2)); !existed {
				node2 = NewNode(id2)
				graph.AddNode(node2)
			}
			err := graph.AddEdge(StringID(id), StringID(id2), weight)
			if err != nil {
				return nil, err
			}
		}
	}
	return graph, nil
}

// Copy returns a copy of graph
//	deep copy
func (g *graph) Copy() Graph {
	nodes := make(map[ID]Node, g.NodeNum())
	edges := make(map[ID]map[ID]float64, g.EdgeNum())
	for k, v := range g.nodes {
		nodes[k] = NewNode(v.ID().String())
	}
	for k, v := range g.edges {
		edge := make(map[ID]float64, len(v))
		for kk, vv := range v {
			edge[kk] = vv
		}
		edges[k] = edge
	}
	return &graph{
		// RWMutex: sync.RWMutex{},
		nodes: nodes,
		edges: edges,
	}
}
