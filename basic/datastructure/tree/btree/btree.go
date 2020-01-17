package btree

import (
	"fmt"
	"godev/basic"
)

// https://en.wikipedia.org/wiki/B-tree

// Item stored int Node
type Item struct {
	Key   interface{}
	Value interface{}
}

func (item *Item) String() string {
	return fmt.Sprintf("{%+v: %+v}", item.Key, item.Value)
}

// Node in BTree
type Node struct {
	Items    []*Item
	Children []*Node
	Parent   *Node
}

// NewNode creates a Node from Item
func NewNode(item *Item) *Node {
	return &Node{
		Items:    []*Item{item},
		Children: []*Node{}, // value ascending from left to right
		Parent:   nil,
	}
}

// BTree data structure
type BTree struct {
	Root       *Node
	Comparator basic.Comparator
	M          int // MaxNumOfChildren on each node, MinNumOfChildren on each non-leaf node = ceil(M / 2), Middle for split = (M - 1) / 2
	ItemsNum   int // Items number of the tree (include all items on all nodes)
}

// NewBTree creates a BTree with M (MaxNumOfChildren on each node >= 3), and a Comparator for node comparision
func NewBTree(M int, comparator basic.Comparator) *BTree {
	if M < 3 {
		panic("Maximum number of children on each node should be at least 3")
	}
	return &BTree{
		Root:       nil,
		Comparator: comparator,
		M:          M,
	}
}

// Empty returns true if no item inside BTree
func (bTree *BTree) Empty() bool {
	return bTree.ItemsNum == 0
}

// Size returns the number of items inside BTree
func (bTree *BTree) Size() int {
	return bTree.ItemsNum
}

func (bTree *BTree) isLeaf(node *Node) bool {
	return len(node.Children) == 0 // single root is also treated as leaf
}

func (bTree *BTree) isFull(node *Node) bool {
	return len(node.Items) == bTree.M-1
}

func (bTree *BTree) requireSplit(node *Node) bool {
	return len(node.Items) > bTree.M-1
}

// Insert inserts one item into BTree, if the root of BTree is nil, the newly inserted item's node will be the root
func (bTree *BTree) Insert(item *Item) {
	if bTree.Root == nil {
		bTree.Root = NewNode(item)
		bTree.ItemsNum++
	} else {
		if bTree.insert(bTree.Root, item) {
			bTree.ItemsNum++
		}
	}
}

func (bTree *BTree) insert(node *Node, item *Item) (inserted bool) {
	if bTree.isLeaf(node) {
		return bTree.insertLeaf(node, item)
	}
	return bTree.insertInternal(node, item)
}

func (bTree *BTree) findPositionByKey(node *Node, key interface{}) (index int, occupied bool) {
	// binary search, fit for larger M
	// for small M, directly traversal may be more efficient
	head, tail := 0, len(node.Items)-1
	mid := 0
	for head <= tail {
		mid = (head + tail) / 2
		switch bTree.Comparator(node.Items[mid].Key, key) {
		case 1:
			tail = mid - 1
		case 0:
			return mid, true
		case -1:
			head = mid + 1
		}
	}
	return head, false
}

func (bTree *BTree) insertItemAtIdxIntoNode(node *Node, index int, item *Item) {
	node.Items = append(node.Items, nil)
	if index < len(node.Items) {
		copy(node.Items[index+1:], node.Items[index:])
	}
	node.Items[index] = item
}

func (bTree *BTree) insertNodeAtIdxIntoChildren(parent *Node, index int, node *Node) {
	parent.Children = append(parent.Children, nil)
	if index < len(parent.Children) {
		copy(parent.Children[index+1:], parent.Children[index:])
	}
	parent.Children[index] = node
}

func (bTree *BTree) insertLeaf(node *Node, item *Item) (inserted bool) {
	insertPos, found := bTree.findPositionByKey(node, item.Key)
	// find item with the same key, just update item value, elements number remains the same
	if found {
		node.Items[insertPos] = item
		return false
	}
	bTree.insertItemAtIdxIntoNode(node, insertPos, item)
	// check split
	bTree.split(node)
	return true
}

func (bTree *BTree) insertInternal(node *Node, item *Item) (inserted bool) {
	insertPos, found := bTree.findPositionByKey(node, item.Key)
	// find item with the same key, just update item value, elements number remains the same
	if found {
		node.Items[insertPos] = item
		return false
	}
	// insert leaf first then split
	return bTree.insert(node.Children[insertPos], item)
}

func (bTree *BTree) split(node *Node) {
	// no need to split
	if !bTree.requireSplit(node) {
		return
	}

	// split parent
	// if parent is root, will create a new root and tree height += 1
	// non-root parent has constrains on its item num:
	//	if M is odd -> mid = (M - 1) / 2
	//	if M is even -> mid = M / 2 - 1 (right sub-tree items number >= left sub-tree items number) = (M - 1) / 2
	//	mid = (M - 1) / 2
	if node == bTree.Root {
		bTree.splitRoot()
	} else {
		bTree.splitNonRoot(node)
	}
}

func (bTree *BTree) splitRoot() {
	mid := (bTree.M - 1) / 2
	// split into left, right sub-trees
	leftSubTree, rightSubTree := bTree.splitIntoLR(bTree.Root, mid)
	// create new root
	newRoot := &Node{
		Items:    []*Item{bTree.Root.Items[mid]},
		Children: []*Node{leftSubTree, rightSubTree},
		Parent:   nil,
	}
	// set sub-trees' parent to new root
	leftSubTree.Parent, rightSubTree.Parent = newRoot, newRoot
	// set new root
	bTree.Root = newRoot
}

func (bTree *BTree) splitIntoLR(node *Node, mid int) (leftSubTree, rightSubTree *Node) {
	// create left, right sub-trees
	leftSubTree = &Node{
		Items:    append([]*Item{}, node.Items[:mid]...), // use append empty slice to ensure result is slice
		Children: nil,
		Parent:   nil,
	}
	rightSubTree = &Node{
		Items:    append([]*Item{}, node.Items[mid+1:]...),
		Children: nil,
		Parent:   nil,
	}
	if len(node.Children) != 0 {
		leftSubTree.Children = append([]*Node{}, node.Children[:mid+1]...)
		rightSubTree.Children = append([]*Node{}, node.Children[mid+1:]...)
		// set subtrees' parent
		setParent(leftSubTree.Children, leftSubTree)
		setParent(rightSubTree.Children, rightSubTree)
	}
	return
}

func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (bTree *BTree) splitNonRoot(node *Node) {
	mid := (bTree.M - 1) / 2
	parent := node.Parent

	// split into left, right sub-trees
	leftSubTree, rightSubTree := bTree.splitIntoLR(node, mid)
	// set sub-trees' parent
	leftSubTree.Parent, rightSubTree.Parent = parent, parent

	// insert node's middle item into parent
	item := node.Items[mid]
	insertPos, _ := bTree.findPositionByKey(parent, item.Key)
	bTree.insertItemAtIdxIntoNode(parent, insertPos, item)

	// set parent's newly inserted item's corresponding node to leftSubTree
	parent.Children[insertPos] = leftSubTree

	// set parent's newly inserted item's next node to rightSubTree
	bTree.insertNodeAtIdxIntoChildren(parent, insertPos+1, rightSubTree)

	// check split
	bTree.split(parent)
}

// Lookup finds whether key inside the tree, it will return the value of the key if found
func (bTree *BTree) Lookup(key interface{}) (value interface{}, found bool) {
	node, index, found := bTree.lookupRec(bTree.Root, key)
	if found {
		return node.Items[index].Value, true
	}
	return nil, false
}

func (bTree *BTree) lookupRec(startNode *Node, key interface{}) (node *Node, index int, found bool) {
	if bTree.Empty() {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = bTree.findPositionByKey(node, key)
		if found {
			return node, index, true
		}
		if bTree.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[index]
	}
}

// Values returns values of all items inside BTree
func (bTree *BTree) Values() []interface{} {
	values := make([]interface{}, bTree.Size())
	var items []*Item
	bTree.items(bTree.Root, &items)
	for idx, i := range items {
		values[idx] = i.Value
	}
	return values
}

func (bTree *BTree) items(node *Node, items *[]*Item) {
	if len(node.Children) != 0 {
		mid := (len(node.Children)-1)/2 + 1
		// left sub-tree
		for _, c := range node.Children[:mid] {
			bTree.items(c, items)
		}
		// mid
		*items = append(*items, node.Items...)
		// right sub-tree
		for _, c := range node.Children[mid:] {
			bTree.items(c, items)
		}
	} else {
		*items = append(*items, node.Items...)
	}
}

// Clear clears items inside BTree by creating a new BTree with the same M and Comparator
func (bTree *BTree) Clear() {
	*bTree = *NewBTree(bTree.M, bTree.Comparator)
}

// Delete delete tree node with key
func (bTree *BTree) Delete(key interface{}) {
	node, index, found := bTree.lookupRec(bTree.Root, key)
	if found {
		bTree.delete(node, index)
		bTree.ItemsNum--
	}
}

func (bTree *BTree) delete(node *Node, index int) {
	// two conditions:
	//	1. delete item from leaf node
	//	2. delete item from internal node
	//	one more step after deletion: re-balancing

	// delete item from leaf node
	if bTree.isLeaf(node) {
		deletedKey := node.Items[index].Key // use deletedItem's key as reference for re-balancing
		bTree.deleteItemAtIdx(node, index)  // delete item
		bTree.rebalance(node, deletedKey)   // re-balance
		// root?
	} else {
		// delete item from internal node
		// 1. choose a new separator (either the largest element in the left subtree or the smallest element in the right subtree),
		// remove it from the leaf node it is in, and replace the element to be deleted with the new separator.
		leftLargestNode := bTree.rightMost(node.Children[index]) // the largest node in the left subtree
		leftLargestItemIndex := len(leftLargestNode.Items) - 1
		leftLargestItem := leftLargestNode.Items[leftLargestItemIndex] // put item in internal node as new separator
		node.Items[index] = leftLargestItem
		deletedKey := leftLargestItem.Key                            // use deletedItem's key as reference for re-balancing
		bTree.deleteItemAtIdx(leftLargestNode, leftLargestItemIndex) // delete that item from leftLargestNode
		bTree.rebalance(node, deletedKey)                            // re-balance
	}
}

func (bTree *BTree) deleteItemAtIdx(node *Node, index int) {
	copy(node.Items[index:], node.Items[index+1:])
	node.Items[len(node.Items)-1] = nil
	node.Items = node.Items[:len(node.Items)-1]
}

func (bTree *BTree) rightMost(node *Node) *Node {
	// recursively
	currentNode := node
	for {
		// rightMost should be leaf
		if bTree.isLeaf(currentNode) {
			return currentNode
		}
		currentNode = currentNode.Children[len(currentNode.Children)-1]
	}
}

func (bTree *BTree) rebalance(node *Node, deletedKey interface{}) {
	// re-balancing starts from a leaf and proceeds toward the root until the tree is balanced.
	// if deleting an element from a node has brought it under the minimum size ([M - 1] / 2),
	// then some elements must be redistributed to bring all nodes up to the minimum.
	minItemNum := (bTree.M - 1) / 2
	if len(node.Items) >= minItemNum {
		return
	}

	// if the deficient node's right sibling exists and has more than the minimum number of elements,
	// then rotate left
	rs, rsIdx := bTree.rightSibling(node, deletedKey)
	if rs != nil && len(rs.Items) > minItemNum {
		// append parent's separator item into node items
		separatorIdx := rsIdx - 1
		node.Items = append(node.Items, node.Parent.Items[separatorIdx])
		node.Parent.Items[separatorIdx] = rs.Items[0] // leftmost item as new separator, move it to parent
		bTree.deleteItemAtIdx(rs, 0)                  // delete it from right sibling
		// move leftmost children of right sibling to parent if it is internal node (leaf node has no children)
		if !bTree.isLeaf(rs) {
			rsLeftMostChild := rs.Children[0]
			rsLeftMostChild.Parent = node
			node.Children = append(node.Children, rsLeftMostChild)
			bTree.deleteChildAtIdx(rs, 0)
		}
		return
	}

	// otherwise, if the deficient node's left sibling exists and has more than the minimum number of elements,
	// then rotate right
	ls, lsIdx := bTree.leftSibling(node, deletedKey)
	if ls != nil && len(ls.Items) > minItemNum {
		// prepend parent's separator item into node items
		separatorIdx := lsIdx
		node.Items = append([]*Item{node.Parent.Items[separatorIdx]}, node.Items...)
		node.Parent.Items[separatorIdx] = ls.Items[len(ls.Items)-1] // rightmost item as new separator, move it to parent
		bTree.deleteItemAtIdx(ls, len(ls.Items)-1)                  // delete it from left sibling
		// move rightmost children of left sibling to parent if it is internal node (leaf node has no children)
		if bTree.isLeaf(ls) {
			lsRightMostChild := ls.Children[len(ls.Children)-1]
			lsRightMostChild.Parent = node
			// prepend
			node.Children = append([]*Node{lsRightMostChild}, node.Children...)
			bTree.deleteChildAtIdx(ls, len(ls.Children)-1)
		}
		return
	}

	// otherwise, if both immediate siblings have only the minimum number of elements,
	// then merge with a sibling sandwiching their separator taken off from their parent
	if rs != nil {
		// merge with right sibling

		separatorIdx := rsIdx - 1

		// deal with item
		// append parent's separator item into node items
		node.Items = append(node.Items, node.Parent.Items[separatorIdx])
		// append right sibling items into node items
		node.Items = append(node.Items, rs.Items...)
		// delete separator item from parent
		deletedKey = node.Parent.Items[separatorIdx].Key
		bTree.deleteItemAtIdx(node.Parent, separatorIdx)

		// deal with child
		// append right sibling's children into node's children
		node.Children = append(node.Children, node.Parent.Children[rsIdx].Children...)
		// set their parent to node
		setParent(node.Parent.Children[rsIdx].Children, node)
		// delete right sibling from parent
		bTree.deleteChildAtIdx(node.Parent, rsIdx)
	} else if ls != nil {
		// merge with left sibling

		separatorIdx := lsIdx

		// deal with item
		// append parent's separator item into left sibling items, introduce a tmp items variable to avoid change on ls.Items directly
		tmpItems := append([]*Item{}, ls.Items...) // ensure result is []*Item
		tmpItems = append(tmpItems, node.Parent.Items[separatorIdx])
		// prepend left sibling items into node items
		node.Items = append(tmpItems, node.Items...)
		// delete separator item from parent
		deletedKey = node.Parent.Items[lsIdx].Key
		bTree.deleteItemAtIdx(node.Parent, lsIdx)

		// deal with child
		// prepend left sibling's children into node's children
		temChildren := append([]*Node{}, node.Parent.Children[lsIdx].Children...) // ensure result is []*Node
		node.Children = append(temChildren, node.Children...)
		// set their parent to node
		setParent(node.Parent.Children[lsIdx].Children, node)
		// delete left sibling from parent
		bTree.deleteChildAtIdx(node.Parent, lsIdx)
	}

	// if the parent is the root and now has no elements,
	// then free it and make the merged node the new root (tree becomes shallower)
	if node.Parent == bTree.Root && len(bTree.Root.Items) == 0 {
		bTree.Root = node
		node.Parent = nil
		return
	}

	// otherwise, if the parent has fewer than the required number of elements,
	// then re-balance the parent
	bTree.rebalance(node.Parent, deletedKey)
}

func (bTree *BTree) rightSibling(node *Node, key interface{}) (sibling *Node, siblingIndex int) {
	if node.Parent != nil {
		index, _ := bTree.findPositionByKey(node, key)
		index++ // right sibling index
		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (bTree *BTree) leftSibling(node *Node, key interface{}) (sibling *Node, siblingIndex int) {
	if node.Parent != nil {
		index, _ := bTree.findPositionByKey(node, key)
		index-- // right sibling index
		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (bTree *BTree) deleteChildAtIdx(node *Node, index int) {
	// check whether under-deleted child is child
	if index >= len(node.Children) {
		return
	}
	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}
