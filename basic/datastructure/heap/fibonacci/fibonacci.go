package fibonacci

import (
	"fmt"
	"goContainer/basic/datastructure/heap"
	"math"
)

// Heap (Fibonacci Heap)
//	https://en.wikipedia.org/wiki/Fibonacci_heap
//	https://www.cnblogs.com/skywang12345/p/3659060.html
//	minimum heap (consist of a series of minimum ordered tree)
type Heap struct {
	// maintain a pointer to the root containing the minimum key
	root    *node
	itemNum int
}

// node consists of Fibonacci Heap, used internally
type node struct {
	item heap.Item
	// children are linked using a circular doubly linked list (parent / child)
	parent, child *node
	// each child has two siblings (left / right)
	left, right *node
	// isMarked indicates whether current node lost first child from last time when it became another node's child
	isMarked bool
	// number of children
	degree int
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

// DeleteMin (Extract Min) delete the root node from doubly linked list and returns its item
//	then find the new root
func (h *Heap) DeleteMin() heap.Item {
	if h.root == nil {
		return nil
	}
	r := h.root
	// add root's children (child and its left / right siblings) to heap's root list
	for {
		if x := r.child; x != nil {
			x.parent = nil
			// sibling
			if x.right != x {
				r.child = x.right
				// cut x from r, then add it to roots list (the following code)
				remove(x)
			} else {
				r.child = nil
			}
			// add x to root doubly linked list (insert)
			insert(x, r)
		} else {
			break
		}
	}

	// remove r from heap's root doubly linked list
	remove(r)

	// check roots after removal
	if r == r.right {
		h.root = nil
	} else {
		h.root = r.right
		// consolidate heap
		h.consolidate()
	}
	h.itemNum--
	return r.item
}

func (h *Heap) consolidate() {
	// use a map to tell whether there exists two sub-trees (roots) which have the same degree
	degreeMap := make(map[int]*node)
	head := h.root
	tail := head.left
	for {
		x := head
		next := head.right
		deg := head.degree
		for {
			if y, found := degreeMap[deg]; !found {
				// there's no roots with the same degree
				break
			} else {
				// the less one is the root (x)
				if y.item.Compare(x.item) < 0 {
					x, y = y, x
				}
				link(y, x)
				delete(degreeMap, deg)
				deg++
			}
		}
		degreeMap[deg] = x
		// loop over the whole list
		if head == tail {
			break
		}
		// move to next
		head = next
	}
	// reconstruct the heap from degreeMap
	h.root = nil
	for _, n := range degreeMap {
		h.insertNode(n)
	}
}

// link node `n` to root `r`
func link(n, r *node) {
	// remove node n from heap's root list
	remove(n)
	// link n to r (n as r's child and increase r's degree)
	n.parent = r
	r.degree++

	if r.child == nil {
		// r has no child, so n is the only child in the children doubly linked list
		r.child = n
		n.left = n
		n.right = n
	} else {
		// insert to r's children doubly linked list
		insert(n, r.child)
	}
	// isMarked indicates whether current node lost first child from last time when it became another node's child
	// here n just becomes r's child, so set n.isMarked to false
	n.isMarked = false
}

// Insert inserts item into heap, just insert it into heap's roots doubly linked list
//	if the item less than min, replace the min
func (h *Heap) Insert(item heap.Item) {
	nNode := &node{
		item:     item,
		parent:   nil,
		child:    nil,
		left:     nil,
		right:    nil,
		isMarked: false,
		degree:   0,
	}
	h.insertNode(nNode)
	h.itemNum++
}

func (h *Heap) insertNode(nNode *node) {
	if h.root == nil {
		nNode.left = nNode
		nNode.right = nNode
		h.root = nNode
		return
	}
	// insert node before root
	insert(nNode, h.root)
	// if the item less than min, replace the min (root)
	if nNode.item.Compare(h.root.item) < 0 {
		h.root = nNode
	}
}

//	insert node before root, which means the `tail` of doubly linked list
func insert(nNode, root *node) {
	nNode.left = root.left
	root.left.right = nNode
	nNode.right = root
	root.left = nNode
}

// cat is different from `insert`, it appends n2 to n1
//	notice: n1 and n2 are doubly linked list node
func cat(n1, n2 *node) {
	var tmp *node
	tmp = n1.right
	n1.right = n2.right
	n2.right.left = n1
	n2.right = tmp
	tmp.left = n2
}

// remove node from its sibling list
func remove(n *node) {
	n.right.left = n.left
	n.left.right = n.right
}

// renew parent node's degree after cutting its child
func renewDegree(parent *node, degree int) {
	parent.degree -= degree
	if parent.parent != nil {
		renewDegree(parent.parent, degree)
	}
}

// DecreaseKey decrease input node's item to nItem
//	1. cut decreased node from its heap and add this node (single node or root of a sub-tree) to roots list
//	2. cascading cut on decreased node's parent node to ensure the min heap property
//	3. update heap root (min)
func (h *Heap) DecreaseKey(n *node, nItem heap.Item) error {
	if n == nil {
		return nil
	}
	if nItem.Compare(n.item) >= 0 {
		return fmt.Errorf("decrease failed: the new item(%+v) is no smaller than current item(%+v)", nItem, n.item)
	}

	n.item = nItem
	if p := n.parent; p != nil && p.item.Compare(n.item) > 0 {
		// cut node from parent node and add node to roots list
		h.cut(n, p)
		h.cascadingCut(p)
	}
	// update heap root (min)
	if n.item.Compare(h.root.item) < 0 {
		h.root = n
	}
	return nil
}

// cut node from its parent and add to the heap's roots list
func (h *Heap) cut(n, p *node) {
	if n == nil || p == nil {
		return
	}
	// remove node from its parent's children list and decrease its parent's degree
	// remove node
	remove(n)
	// renew degree
	renewDegree(p, n.degree)
	if n.right == n {
		// no sibling
		p.child = nil
	} else {
		p.child = n.right
	}

	n.parent = nil
	n.left, n.right = n, n
	n.isMarked = false

	// add n to roots list
	insert(n, h.root)
}

// cascadingCut recursively cut nodes starting from the root of tree whose child has been cut
func (h *Heap) cascadingCut(n *node) {
	if n == nil {
		return
	}

	p := n.parent

	if p != nil {
		return
	}

	if n.isMarked == false {
		// n has been cut a child
		n.isMarked = true
	} else {
		h.cut(n, p)
		h.cascadingCut(p)
	}
}

// IncreaseKey decrease input node's item to nItem
//	1. add increased node's children (child and child's siblings) into roots list
//	2. cut (cut and cascadingCut) increased node and add it into roots list
//	3. update heap root (min) if n is root
func (h *Heap) IncreaseKey(n *node, nItem heap.Item) error {
	if nItem.Compare(n.item) <= 0 {
		return fmt.Errorf("increase failed: the new item(%+v) is no larger than current item(%+v)", nItem, n.item)
	}

	for {
		child := n.child

		if child == nil {
			break
		}

		// remove child from children list
		remove(child)

		// update n.child
		if child.right == child {
			n.child = nil
		} else {
			n.child = child.right
		}

		// add child into roots list
		insert(child, h.root)
		child.parent = nil
	}

	// add node into roots list
	n.degree = 0
	n.item = nItem
	if p := n.parent; p != nil {
		h.cut(n, p)
		h.cascadingCut(p)
	} else {
		// update heap root (min) if n is root
		if h.root == n {
			right := n.right
			for right != n {
				if h.root.item.Compare(right.item) > 0 {
					h.root = right
				}
				right = right.right
			}
		}
	}
	return nil
}

// Update returns true if `item` is successfully replaced by `nItem`
func (h *Heap) Update(item, nItem heap.Item) bool {
	n := h.search(h.root, item)
	if n == nil {
		fmt.Printf("input node is nil and item is %+v", nItem)
		return false
	}
	if nItem.Compare(n.item) < 0 {
		if err := h.DecreaseKey(n, nItem); err != nil {
			fmt.Println(err)
			return false
		}
		return true
	} else if nItem.Compare(n.item) > 0 {
		if err := h.IncreaseKey(n, nItem); err != nil {
			fmt.Println(err)
			return false
		}
		return true
	} else {
		fmt.Println("no need to update")
		return false
	}
}

// Meld returns union of two heaps (notice: in place change, input heaps may be changed)
//	for efficiency consideration, add to heap which has larger maxDegree to achieve less operations
//	TODO: copy better than in-place change?
func (h *Heap) Meld(ah *Heap) *Heap {
	if h.root == nil {
		return ah
	}
	if ah.root == nil {
		return h
	}

	h1, h2 := h, ah

	if ah.maxDegree() > h.maxDegree() {
		h1, h2 = h2, h1
	}

	cat(h2.root, h1.root)
	if h2.root.item.Compare(h1.root.item) < 0 {
		h1.root = h2.root
	}
	h1.itemNum += h2.itemNum
	*h = *h1
	return h
}

// maxDegree estimates max degree of heap
func (h *Heap) maxDegree() int {
	return int(math.Log2(float64(h.itemNum))) + 1 // ceil
}

// Search returns true if input item is found
//	recursively
func (h *Heap) Search(item heap.Item) bool {
	if n := h.search(h.root, item); n != nil {
		return true
	}
	return false
}

// search item in sub-tree
//	returns node which contains the item or nil
func (h *Heap) search(r *node, item heap.Item) *node {
	if r == nil || r.item == item {
		return r
	}

	n := r
	var p *node

	// search in the doubly linked list
	for {
		if n.item == item {
			p = n
			break
		} else {
			// search in sub-sub-tree
			if p = h.search(n.child, item); p != nil {
				break
			}
		}
		n = n.right
		if n == r {
			break
		}
	}
	return p
}

// Delete deletes item from heap and return it
//	need input an additional itemMinimum (minimum value of item type) to assist deletion
func (h *Heap) Delete(item, itemMinimum heap.Item) heap.Item {
	if n := h.search(h.root, item); n != nil {
		return h.delete(n, item, itemMinimum)
	}
	// item not stored inside the heap
	return item
}

// delete deletes node from heap and returns its item
//	1. decrease node's item to a value less than min (root's item)
//	2. call DeleteMin
// this function need to define a itemMinimum
func (h *Heap) delete(n *node, item, itemMinimum heap.Item) heap.Item {
	err := h.DecreaseKey(n, itemMinimum)
	if err != nil {
		return nil
	}
	h.DeleteMin()
	return item
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

// PopAllItems pops all items out the heap with ascending order
//	Notice: the heap will be cleared after calling this method
func (h *Heap) PopAllItems() []heap.Item {
	num := h.itemNum
	var res []heap.Item
	for i := 0; i < num; i++ {
		res = append(res, h.DeleteMin())
	}
	return res
}

// traverse sub-tree
//	store nodes into array
func (h *Heap) traverse(startNode *node, nodes *[]*node) {
	if startNode == nil {
		return
	}
	n := startNode
	for {
		*nodes = append(*nodes, n)
		// traverse in its sub-tree
		if x := n.child; x != nil {
			h.traverse(x, nodes)
		}
		// sibling
		n = n.right
		if n == startNode {
			break
		}
	}
}

// Print prints heap items
func (h *Heap) Print() {
	fmt.Println("***  Fibonacci Heap  ***")
	if h.root == nil {
		fmt.Println("empty heap")
		return
	}

	// print roots list
	r := h.root
	for {
		fmt.Println()
		fmt.Printf("%+v(%+v) is heap's root\n", r.item, r.degree)
		h.print(r.child, r, 0)

		r = r.right
		if r == h.root {
			break
		}
	}
	fmt.Println()
}

// print prints node sub-tree
// direction: child (0) / sibling (1)
func (h *Heap) print(startNode, prevNode *node, direction int) {
	if startNode == nil {
		return
	}
	n := startNode
	for {
		if direction == 0 {
			fmt.Printf("%+v(%+v) is %+v's child\n", n.item, n.degree, prevNode.item)
		} else {
			fmt.Printf("%+v(%+v) is %+v's next\n", n.item, n.degree, prevNode.item)
		}

		// print sub-tree
		if child := n.child; child != nil {
			h.print(child, n, 0)
		}

		// sibling
		prevNode = n
		n = n.right
		direction = 1

		if n == startNode {
			break
		}
	}
}
