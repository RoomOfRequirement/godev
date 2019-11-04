package rbtree

import (
	"fmt"
	"goContainer"
)

// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree

// Red, Black color of Red-Black Tree, represented by bool
const (
	Red   = false
	Black = true
)

// Node struct
type Node struct {
	key                         interface{}
	value                       interface{}
	color                       bool
	leftTree, rightTree, parent *Node
}

// NIL represents nil Node, which used for finding leaf
var NIL = &Node{
	key:       nil,
	value:     nil,
	color:     Black,
	leftTree:  nil,
	rightTree: nil,
	parent:    nil,
}

// NewNode creates a new Node from key
func NewNode(key, value interface{}) *Node {
	return &Node{
		key:       key,
		value:     value,
		color:     Red,
		leftTree:  nil,
		rightTree: nil,
		parent:    nil,
	}
}

func (node *Node) grandparent() *Node {
	if node.parent == nil {
		return nil
	}
	return node.parent.parent
}

func (node *Node) uncle() *Node {
	if node.grandparent() == nil {
		return nil
	}
	if node.parent == node.grandparent().rightTree {
		return node.grandparent().leftTree
	}
	return node.grandparent().rightTree
}

func (node *Node) sibling() *Node {
	if node.parent.leftTree == node {
		return node.parent.rightTree
	}
	return node.parent.leftTree
}

// String returns node list in-order
func (node *Node) String() string {
	if node == NIL {
		return "NIL"
	}
	s := ""
	if node.leftTree != nil {
		s += node.leftTree.String() + " "
	}
	c := "Red"
	if node.color {
		c = "Black"
	}
	s += fmt.Sprintf("%v(%s)", node.key, c)
	if node.rightTree != nil {
		s += " " + node.rightTree.String()
	}
	return "[" + s + "]"
}

// RBTree struct stands for Red-Black Tree data structure
type RBTree struct {
	Root       *Node
	Comparator container.Comparator
}

// String returns nodes in rbTree, not pretty, need improved
func (rbTree *RBTree) String() string {
	return fmt.Sprintf("%s", rbTree.Root)
}

func (rbTree *RBTree) value(node *Node, dataSlice *[]interface{}) {
	if node.key == nil {
		return
	}
	rbTree.value(node.leftTree, dataSlice)
	*dataSlice = append(*dataSlice, node.value)
	rbTree.value(node.rightTree, dataSlice)
}

// Values returns values of all nodes inside the tree
//	notice: values follows keys' order!
func (rbTree *RBTree) Values() []interface{} {
	if rbTree.Root == nil {
		return nil
	}
	var values []interface{}
	rbTree.value(rbTree.Root, &values)
	return values
}

func (rbTree *RBTree) key(node *Node, dataSlice *[]interface{}) {
	if node.key == nil {
		return
	}
	rbTree.key(node.leftTree, dataSlice)
	*dataSlice = append(*dataSlice, node.key)
	rbTree.key(node.rightTree, dataSlice)
}

// Keys returns keys of all nodes inside the tree
func (rbTree *RBTree) Keys() []interface{} {
	if rbTree.Root == nil {
		return nil
	}
	var keys []interface{}
	rbTree.key(rbTree.Root, &keys)
	return keys
}

func (rbTree *RBTree) size(node *Node) int {
	if node == nil || node == NIL {
		return 0
	}
	return rbTree.size(node.leftTree) + 1 + rbTree.size(node.rightTree)
}

// Size returns number of nodes inside the tree
func (rbTree *RBTree) Size() int {
	return rbTree.size(rbTree.Root)
}

// Empty returns true if the tree has no nodes
func (rbTree *RBTree) Empty() bool {
	return rbTree.Root == nil
}

func (rbTree *RBTree) minKey(node *Node) interface{} {
	if node.leftTree == NIL {
		return node.key
	}
	return rbTree.minKey(node.leftTree)
}

// MinKey returns the minimum key inside nodes of the tree
func (rbTree *RBTree) MinKey() interface{} {
	return rbTree.minKey(rbTree.Root)
}

// NewRBTree creates a new red-black tree
func NewRBTree(root *Node, comparator container.Comparator) *RBTree {
	rbTree := &RBTree{Root: root, Comparator: comparator}
	rbTree.Root.color = Black
	rbTree.Root.leftTree, rbTree.Root.rightTree = NIL, NIL
	return rbTree
}

/*
Rotate left on Y:
     gp                gp
     /                 /
    X                 Y
   / \               / \
  a   Y    ----->   X   c
     / \           / \
    b   c         a   b
*/
func (rbTree *RBTree) rotateLeft(Y *Node) {
	if Y.parent == nil {
		rbTree.Root = Y
		return
	}

	gp := Y.grandparent()
	X := Y.parent
	b := Y.leftTree

	X.rightTree = b
	if b != NIL {
		b.parent = X
	}
	Y.leftTree = X
	X.parent = Y

	if rbTree.Root == X {
		rbTree.Root = Y
	}
	Y.parent = gp

	if gp != nil {
		if gp.leftTree == X {
			gp.leftTree = Y
		} else {
			gp.rightTree = Y
		}
	}
}

/*
Rotate right on Y:
     gp                gp
     /                 /
    X                 Y
   / \               / \
  Y   a    ----->   b   X
 / \                   / \
b   c                 c   a
*/
func (rbTree *RBTree) rotateRight(Y *Node) {
	gp := Y.grandparent()
	X := Y.parent
	c := Y.rightTree

	X.leftTree = c

	if c != NIL {
		c.parent = X
	}
	Y.rightTree = X
	X.parent = Y

	if rbTree.Root == X {
		rbTree.Root = Y
	}
	Y.parent = gp

	if gp != nil {
		if gp.leftTree == X {
			gp.leftTree = Y
		} else {
			gp.rightTree = Y
		}
	}
}

/*
Y color: Red
     gp                gp
     /                 /
    X                 X
   / \               / \
  a   Y    <---->   Y   a
     / \           / \
    b   c         b   c
*/
func (rbTree *RBTree) insertCase(Y *Node) {
	if Y.parent == nil {
		rbTree.Root = Y
		Y.color = Black
		return
	}

	gp := Y.grandparent()
	X := Y.parent
	a := Y.uncle()
	b := Y.leftTree
	c := Y.rightTree

	if X.color == Red {
		if a.color == Red {
			X.color, a.color = Black, Black
			gp.color = Red
			rbTree.insertCase(gp)
		} else {
			// rotate to left
			if X.rightTree == Y && gp.leftTree == X {
				rbTree.rotateLeft(Y)
				Y.color = Black
				b.color, c.color = Red, Red
			} else if X.leftTree == Y && gp.rightTree == X {
				// rotate to right
				rbTree.rotateRight(Y)
				Y.color = Black
				b.color, c.color = Red, Red
			} else if X.leftTree == Y && gp.leftTree == X {
				X.color = Black
				gp.color = Red
				rbTree.rotateRight(X)
			} else if X.rightTree == Y && gp.rightTree == X {
				X.color = Black
				gp.color = Red
				rbTree.rotateLeft(X)
			}
		}
	}
}

func (rbTree *RBTree) insert(node *Node, key, value interface{}) {
	if rbTree.Comparator(node.key, key) >= 0 {
		if node.leftTree != NIL {
			rbTree.insert(node.leftTree, key, value)
		} else {
			tmp := NewNode(key, value)
			tmp.leftTree, tmp.rightTree = NIL, NIL
			tmp.parent = node
			node.leftTree = tmp
			rbTree.insertCase(tmp)
		}
	} else {
		if node.rightTree != NIL {
			rbTree.insert(node.rightTree, key, value)
		} else {
			tmp := NewNode(key, value)
			tmp.leftTree, tmp.rightTree = NIL, NIL
			tmp.parent = node
			node.rightTree = tmp
			rbTree.insertCase(tmp)
		}
	}
}

// Insert inserts key inside the RBTree
func (rbTree *RBTree) Insert(key, value interface{}) {
	if rbTree.Root == nil {
		rbTree.Root = NewNode(key, value)
		rbTree.Root.color = Black
		rbTree.Root.leftTree, rbTree.Root.rightTree = NIL, NIL
	}
	rbTree.insert(rbTree.Root, key, value)
}

func (rbTree *RBTree) getSmallestChild(root *Node) *Node {
	if root.leftTree == NIL {
		return root
	}
	return rbTree.getSmallestChild(root.leftTree)
}

func (rbTree *RBTree) deleteChild(node *Node, key interface{}) bool {
	if rbTree.Comparator(node.key, key) > 0 {
		if node.leftTree == NIL {
			return false
		}
		return rbTree.deleteChild(node.leftTree, key)
	} else if rbTree.Comparator(node.key, key) < 0 {
		if node.rightTree == NIL {
			return false
		}
		return rbTree.deleteChild(node.rightTree, key)
	} else if rbTree.Comparator(node.key, key) == 0 {
		if node.rightTree == NIL {
			rbTree.deleteOneChild(node)
			return true
		}
		smallestNode := rbTree.getSmallestChild(node.rightTree)
		smallestNode.key, node.key = node.key, smallestNode.key
		rbTree.deleteOneChild(smallestNode)
		return true
	}
	return false
}

// when Y has at most one non-NIL child
func (rbTree *RBTree) deleteOneChild(Y *Node) {
	Child := NIL
	if Y.leftTree == NIL {
		Child = Y.rightTree
	} else {
		Child = Y.leftTree
	}

	X := Y.parent

	if X == nil && Y.leftTree == NIL && Y.rightTree == NIL {
		Y = nil
		rbTree.Root = Y
		return
	}

	if X == nil {
		Y = nil
		Child.parent = nil
		rbTree.Root = Child
		rbTree.Root.color = Black
		return
	}

	if X.leftTree == Y {
		X.leftTree = Child
	} else {
		X.rightTree = Child
	}
	Child.parent = X

	if Y.color == Black {
		if Child.color == Red {
			Child.color = Black
		} else {
			rbTree.deleteCase(Child)
		}
	}
	Y = nil
}

func (rbTree *RBTree) deleteCase(Y *Node) {
	X := Y.parent
	if X == nil {
		Y.color = Black
		return
	}

	S := Y.sibling()

	if S.color == Red {
		X.color = Red
		S.color = Black
		if Y == X.leftTree {
			rbTree.rotateLeft(X)
		} else {
			rbTree.rotateRight(X)
		}
	}

	if S != NIL {
		if X.color == Black && S.color == Black && S.leftTree.color == Black && S.rightTree.color == Black {
			S.color = Red
			rbTree.deleteCase(X)
		} else if X.color == Red && S.color == Black && S.leftTree.color == Black && S.rightTree.color == Black {
			S.color = Red
			X.color = Black
		} else {
			if S.color == Black {
				if Y == X.leftTree && S.leftTree.color == Red && S.rightTree.color == Black {
					S.color = Red
					S.leftTree.color = Black
					rbTree.rotateRight(S.leftTree)
				} else if Y == X.rightTree && S.leftTree.color == Black && S.rightTree.color == Red {
					S.color = Red
					S.rightTree.color = Black
					rbTree.rotateLeft(S.rightTree)
				}
			}

			S.color = X.color
			X.color = Black
			if Y == X.leftTree {
				S.rightTree.color = Black
				rbTree.rotateLeft(S)
			} else {
				S.leftTree.color = Black
				rbTree.rotateRight(S)
			}
		}
	} else {
		if X.color == Black {
			rbTree.deleteCase(X)
		} else if X.color == Red {
			X.color = Black
		}
	}
}

// Delete returns true if the input key inside the tree's nodes and successfully deleted
func (rbTree *RBTree) Delete(key interface{}) bool {
	return rbTree.deleteChild(rbTree.Root, key)
}

// Clear clears all nodes inside the tree by setting root to nil
func (rbTree *RBTree) Clear() {
	rbTree.Root = nil
}
