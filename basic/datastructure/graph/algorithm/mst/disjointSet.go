package mst

import (
	"goContainer/basic/datastructure/graph"
)

// DisjointSet struct
//	https://en.wikipedia.org/wiki/Disjoint-set_data_structure
type DisjointSet struct {
	parent  *DisjointSet
	rank    int
	Element map[graph.ID]struct{} // can be anything, here i use a map of graph.ID for realizing Kruskal
}

// MakeSet operation makes a new set by creating a new element with a unique id, a rank of 0, and a parent pointer to itself
// 	the parent pointer to itself indicates that the element is the representative member of its own set
func MakeSet() *DisjointSet {
	e := &DisjointSet{
		parent:  nil,
		rank:    0,
		Element: make(map[graph.ID]struct{}),
	}
	e.parent = e
	return e
}

// Find follows the chain of parent pointers from x up the tree until it reaches a root element, whose parent is itself
//	use path splitting
func (ds *DisjointSet) Find() *DisjointSet {
	for ds.parent != ds {
		ds, ds.parent = ds.parent, ds.parent.parent
	}
	return ds
}

// Union uses Find to determine the roots of the trees x and y belong to
//	if the roots are distinct, the trees are combined by attaching the root of one to the root of the other
//	by rank
func Union(ds1, ds2 *DisjointSet) {
	ds1Root, ds2Root := ds1.Find(), ds2.Find()

	// ds1 and ds2 are already in the same set
	if ds1Root == ds2Root {
		return
	}

	// ds1 and ds2 are not in the same set, so merge them
	// merge set with smaller rank into larger oen
	if ds1Root.rank < ds2Root.rank {
		ds1Root, ds2Root = ds2Root, ds1Root
	}
	ds2Root.parent = ds1Root
	if ds1Root.rank == ds2Root.rank {
		ds1Root.rank++
	}
}

// DSSet is a set of DisjointSet
type DSSet struct {
	Sets map[*DisjointSet]struct{}
}

// NewDSSet returns a new set of disjoint set
func NewDSSet() *DSSet {
	return &DSSet{Sets: make(map[*DisjointSet]struct{})}
}

// FindSet returns the disjoint set which contains input ID
func (dss *DSSet) FindSet(id graph.ID) *DisjointSet {
	for ds := range dss.Sets {
		if _, found := ds.Element[id]; found {
			return ds
		}
	}
	return nil
}
