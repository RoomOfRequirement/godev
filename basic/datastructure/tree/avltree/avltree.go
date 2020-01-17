package avltree

import (
	"godev/basic"
)

// AVLTree struct
//	https://en.wikipedia.org/wiki/AVL_tree
type AVLTree struct {
	Root       *node
	Comparator basic.Comparator
	itemNum    int
}

type node struct {
	key, value                  interface{}
	leftTree, rightTree, parent *node
	height                      int
}

func newNode(key, value interface{}) *node {
	return &node{
		key:       key,
		value:     value,
		leftTree:  nil,
		rightTree: nil,
		parent:    nil,
		height:    1, // leaf
	}
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func getHeight(n *node) int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *node) getBalance() int {
	return getHeight(n.leftTree) - getHeight(n.rightTree)
}

func (n *node) updateHeight() {
	n.height = maxInt(getHeight(n.leftTree), getHeight(n.rightTree)) + 1
}

/*
     n                y
    / \              / \
   a   y   ----->   n   c
      / \          / \
     b   c        a   b
*/
func (n *node) leftRotate() *node {
	y := n.rightTree
	b := y.leftTree

	y.leftTree = n
	n.rightTree = b

	n.updateHeight()
	y.updateHeight()

	return y
}

/*
     n                n
    / \              / \
   x   c   ----->   a   y
  / \                  / \
 a   b                b   c
*/
func (n *node) rightRotate() *node {
	x := n.leftTree
	b := x.rightTree

	x.rightTree = n
	n.leftTree = b

	x.updateHeight()
	n.updateHeight()

	return x
}

// NewAVLTree creates a new AVL tree
func NewAVLTree(comparator basic.Comparator) *AVLTree {
	return &AVLTree{
		Root:       nil,
		Comparator: comparator,
		itemNum:    0,
	}
}

// Get returns value with input key when found key inside the tree
func (avlTree *AVLTree) Get(key interface{}) (value interface{}, found bool) {
	node := avlTree.get(key)
	if node != nil {
		return node.value, true
	}
	return nil, false
}

func (avlTree *AVLTree) get(key interface{}) *node {
	node := avlTree.Root
	for node != nil {
		switch avlTree.Comparator(node.key, key) {
		case 0:
			return node
		case -1:
			node = node.rightTree
		case 1:
			node = node.leftTree
		}
	}
	return nil
}

// Set sets node's value to input value of input key, if key exists, change its value to new value
func (avlTree *AVLTree) Set(key, value interface{}) {
	avlTree.Root = avlTree.set(avlTree.Root, key, value)
}

func (avlTree *AVLTree) set(n *node, key, value interface{}) *node {
	if n == nil {
		avlTree.itemNum++
		return newNode(key, value)
	}

	// normal BST insert
	switch avlTree.Comparator(key, n.key) {
	case -1:
		n.leftTree = avlTree.set(n.leftTree, key, value)
	case 1:
		n.rightTree = avlTree.set(n.rightTree, key, value)
	case 0:
		// equal keys, update value
		n.value = value
		return n
	}

	// update height of this ancestor node
	n.updateHeight()

	// use balance factor to see whether this node became unbalanced
	// left left
	if n.getBalance() > 1 && avlTree.Comparator(key, n.leftTree.key) == -1 {
		return n.rightRotate()
	}
	// right right
	if n.getBalance() < -1 && avlTree.Comparator(key, n.rightTree.key) == 1 {
		return n.leftRotate()
	}
	// left right
	if n.getBalance() > 1 && avlTree.Comparator(key, n.leftTree.key) == 1 {
		n.leftTree.leftRotate()
		return n.rightRotate()
	}
	// right left
	if n.getBalance() < -1 && avlTree.Comparator(key, n.rightTree.key) == -1 {
		n.rightTree.rightRotate()
		return n.leftRotate()
	}

	// balanced -> unchanged
	return n
}

// Delete deletes k,v pair if found inside the tree
func (avlTree *AVLTree) Delete(key interface{}) bool {
	if avlTree.Root == nil {
		return false
	}
	avlTree.Root = avlTree.delete(avlTree.Root, key)
	avlTree.itemNum--
	return true
}

func (avlTree *AVLTree) delete(n *node, key interface{}) *node {
	if n == nil {
		return nil
	}
	switch avlTree.Comparator(key, n.key) {
	case -1:
		n.leftTree = avlTree.delete(n.leftTree, key)
	case 1:
		n.rightTree = avlTree.delete(n.rightTree, key)
	case 0:
		// normal BST delete
		// one child or no child
		if n.leftTree == nil && n.rightTree == nil {
			return nil
		} else if n.leftTree == nil {
			n = n.rightTree
		} else if n.rightTree == nil {
			n = n.leftTree
		} else {
			// two children
			leftMost := avlTree.leftMost(n.rightTree)
			n.key = leftMost.key
			n.rightTree = avlTree.delete(n.rightTree, leftMost.key)
		}
	}

	// update height
	n.updateHeight()

	// use balance factor to see whether this node became unbalanced
	// left left
	if n.getBalance() > 1 && n.leftTree.getBalance() >= 0 {
		return n.rightRotate()
	}
	// right right
	if n.getBalance() < -1 && n.rightTree.getBalance() <= 0 {
		return n.leftRotate()
	}
	// left right
	if n.getBalance() > 1 && n.leftTree.getBalance() < 0 {
		n.leftTree.leftRotate()
		return n.rightRotate()
	}
	// right left
	if n.getBalance() < -1 && n.rightTree.getBalance() > 0 {
		n.rightTree.rightRotate()
		return n.leftRotate()
	}

	return n
}

// left most (min) in sub-trees
func (avlTree *AVLTree) leftMost(n *node) *node {
	current := n
	for current.leftTree != nil {
		current = current.leftTree
	}
	return current
}

// Empty returns true if the tree has no k, v pair inside
func (avlTree *AVLTree) Empty() bool {
	return avlTree.itemNum == 0
}

// Size returns number of k, v pairs inside the tree
func (avlTree *AVLTree) Size() int {
	return avlTree.itemNum
}

// Clear clears the tree
func (avlTree *AVLTree) Clear() {
	*avlTree = *NewAVLTree(avlTree.Comparator)
}

// Values returns values of all nodes inside the tree
//	notice: values follows keys' order!
func (avlTree *AVLTree) Values() []interface{} {
	if avlTree.Root == nil {
		return nil
	}
	var values []interface{}
	avlTree.value(avlTree.Root, &values)
	return values
}

func (avlTree *AVLTree) value(node *node, dataSlice *[]interface{}) {
	if node == nil {
		return
	}
	avlTree.value(node.leftTree, dataSlice)
	*dataSlice = append(*dataSlice, node.value)
	avlTree.value(node.rightTree, dataSlice)
}

// Keys returns keys of all nodes inside the tree
func (avlTree *AVLTree) Keys() []interface{} {
	if avlTree.Root == nil {
		return nil
	}
	var keys []interface{}
	avlTree.key(avlTree.Root, &keys)
	return keys
}

func (avlTree *AVLTree) key(node *node, dataSlice *[]interface{}) {
	if node == nil {
		return
	}
	avlTree.key(node.leftTree, dataSlice)
	*dataSlice = append(*dataSlice, node.key)
	avlTree.key(node.rightTree, dataSlice)
}
