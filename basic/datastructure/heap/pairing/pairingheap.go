package pairing

import (
	"goContainer/basic/datastructure/heap"
)

// node consists of Pairing Heap, used internally
//	here i implement it in a way what wiki presents
//	another way to implement is 3 points (lefChild, nextSibling, prev(point to parent)) or 2 pointers (lefChild, nextSibling)
//	but it is not so straightforward to understand
type node struct {
	item heap.Item
	// children list (nodes expect heap top)
	children []*node
	// parent heap
	parent *node
}

// Heap (Pairing Heap) struct
//	https://en.wikipedia.org/wiki/Pairing_heap
//	heap consists of root and sub-heaps
type Heap struct {
	root    *node
	itemNum int
}

// NewHeap creates a new empty heap
func NewHeap() *Heap {
	return &Heap{
		root:    nil,
		itemNum: 0,
	}
}

// FindMin returns the minimum item inside the heap (root)
func (h *Heap) FindMin() heap.Item {
	if h.root == nil {
		return nil
	}
	return h.root.item
}

// Meld returns union of two heaps (notice: in place change, input heaps may be changed)
func (h *Heap) Meld(h1 *Heap) *Heap {
	if h.root == nil {
		return h1
	}
	if h1.root == nil {
		return h
	}
	itemNum := h.itemNum + h1.itemNum
	h.root = meld(h.root, h1.root)
	h.itemNum = itemNum
	return h
}

func meld(r1, r2 *node) *node {
	// if root is nil
	if r1 == nil {
		r1 = r2
		return r2
	}

	if r1.item.Compare(r2.item) < 0 {
		// put r2 as the first child of r1
		r1.children = append([]*node{r2}, r1.children...)
		// update melded heap's parent
		r2.parent = r1
		return r1
	}
	// put r1 as the first child of r2
	r2.children = append([]*node{r1}, r2.children...)
	// update melded heap's parent
	r1.parent = r2
	return r2
}

// Insert inserts an item into heap is by meld the heap with a new heap containing just this item
func (h *Heap) Insert(item heap.Item) {
	h.root = meld(h.root, &node{
		item:     item,
		children: nil,
		parent:   nil,
	})
	h.itemNum++
}

// DeleteMin (Extract Min) delete the root node and update heap structure
//	this requires performing repeated melds of its children until only one tree remains
//	the standard strategy first melds the subheaps in pairs (this is the step that gave this data structure its name) from left to right
//	and then melds the resulting list of heaps from right to left
func (h *Heap) DeleteMin() heap.Item {
	if h.root == nil {
		return nil
	}
	item := h.root.item
	h.root = mergePairs(h.root.children)
	h.itemNum--
	return item
}

func mergePairs(subHeaps []*node) *node {
	if len(subHeaps) == 0 {
		return nil
	} else if len(subHeaps) == 1 {
		subHeaps[0].parent = nil
		return subHeaps[0]
	} else {
		var merged *node
		for {
			// meld(meld(list[0], list[1]), merge-pairs(list[2..]))
			if len(subHeaps) == 0 {
				break
			}
			if merged == nil {
				merged = meld(subHeaps[0], subHeaps[1])
				subHeaps = subHeaps[2:]
			} else {
				merged = meld(merged, subHeaps[0])
				subHeaps = subHeaps[1:]
			}
		}
		merged.parent = nil
		return merged
	}
}

// Search returns true if input item is found
//	recursively
func (h *Heap) Search(item heap.Item) bool {
	if n := h.root.search(item); n != nil {
		return true
	}
	return false
}

// search item in sub-heaps
//	returns node which contains the item or nil
func (n *node) search(item heap.Item) *node {
	if n.item.Compare(item) == 0 {
		return n
	}

	if len(n.children) == 0 {
		return nil
	}

	var node *node

loop:
	for _, child := range n.children {
		node = child.search(item)
		if node != nil {
			break loop
		}
	}
	return node
}

// Delete deletes item from heap and return it
func (h *Heap) Delete(item heap.Item) heap.Item {
	node := h.root.search(item)
	if node == nil {
		return nil
	}
	// new children list
	children := node.cut()
	// add to root
	h.root.children = append(h.root.children, children...)
	h.itemNum--
	// set root to nil if empty
	if h.itemNum == 0 {
		h.root = nil
	}
	return node.item
}

// cut node and return its children list
func (n *node) cut() []*node {
	// avoid cut root
	if n.parent == nil {
		return nil
	}
	// update children's parent
	for _, node := range n.children {
		node.parent = nil
	}
	var idx int
	for i, node := range n.parent.children {
		if node == n {
			idx = i
			break
		}
	}
	// cut node from parent's children list
	n.parent.children = append(n.parent.children[:idx], n.parent.children[idx+1:]...)
	n.parent = nil
	// return updated children list
	return n.children
}

// Update returns true if `item` is successfully replaced by `nItem`
func (h *Heap) Update(item, nItem heap.Item) bool {
	n := h.root.search(item)
	if n == nil {
		return false
	}

	if n == h.root {
		h.DeleteMin()
		h.Insert(nItem)
	} else {
		children := n.cut()
		h.Insert(nItem)
		h.root.children = append(h.root.children, children...)
	}
	return true
}

// Size returns item number inside the heap
func (h *Heap) Size() int {
	return h.itemNum
}

// Empty returns true if no item inside heap
func (h *Heap) Empty() bool {
	return h.root == nil
}

// Clear clears heap by setting its root to nil
func (h *Heap) Clear() {
	h.root = nil
	h.itemNum = 0
}

// Values returns values inside the heap
func (h *Heap) Values() []interface{} {
	nodes := make([]*node, 0, h.itemNum)
	h.traverse(h.root, &nodes)
	values := make([]interface{}, h.itemNum)
	for i, n := range nodes {
		values[i] = n.item
	}
	return values
}

// traverse sub-tree
//	store nodes into array
func (h *Heap) traverse(startNode *node, nodes *[]*node) {
	if startNode == nil {
		return
	}
	n := startNode
	*nodes = append(*nodes, n)

	// traverse in its sub-heaps
	n.traverseChildren(nodes)
}

func (n *node) traverseChildren(nodes *[]*node) {
	if n.children == nil || len(n.children) == 0 {
		return
	}
	*nodes = append(*nodes, n.children...)
	for _, child := range n.children {
		child.traverseChildren(nodes)
	}
}
