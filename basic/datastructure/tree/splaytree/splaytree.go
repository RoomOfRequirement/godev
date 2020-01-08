package splaytree

import (
	"fmt"
	"goContainer/basic"
)

// references:
// https://en.wikipedia.org/wiki/Splay_tree
// https://www.geeksforgeeks.org/splay-tree-set-1-insert/

// KeyAlreadyExistError error
type KeyAlreadyExistError struct {
	key interface{}
}

// Error to meet error interface
func (e KeyAlreadyExistError) Error() string {
	return fmt.Sprintf("insert error: key %v already exists, you can choose Update", e.key)
}

// KeyNotExistError error
type KeyNotExistError struct {
	key interface{}
}

// Error to meet error interface
func (e KeyNotExistError) Error() string {
	return fmt.Sprintf("delete error: key %v does NOT exist", e.key)
}

// node struct used inside tree
type node struct {
	key, value          interface{}
	left, right, parent *node
}

func newNode(k, v interface{}) *node {
	return &node{
		key:    k,
		value:  v,
		left:   nil,
		right:  nil,
		parent: nil,
	}
}

// SplayTree struct
//	TODO: all containers need comparision operations, where to bind the comparision functions?
//	 better on Node(Item) or Tree(Container)?
type SplayTree struct {
	root       *node
	comparator basic.Comparator
	size       int
}

// NewSplayTree creates a new empty splay tree
func NewSplayTree(comparator basic.Comparator) *SplayTree {
	return &SplayTree{
		root:       nil,
		comparator: comparator,
		size:       0,
	}
}

/*
p (x.parent) is root
zig: left case (x is left child)
right rotate (p):
    p                 x
   / \               / \
  x   c    ----->   a   p
 / \                   / \
a   b                 b   c
*/
func (st *SplayTree) zig(x *node) {
	x.parent.left = x.right // p.left = b
	if x.right != nil {     // if b != nil
		x.right.parent = x.parent // b.parent = p
	}
	x.parent.parent = x // p.parent = x
	x.right = x.parent  // x.right = p

	x.parent = nil
	st.root = x
}

/*
p (x.parent) is root
zag: right case (x is right child)
left rotate (p):
    p                 x
   / \               / \
  c   x    ----->   p   a
     / \           / \
    a   b         c   b
*/
func (st *SplayTree) zag(x *node) {
	x.parent.right = x.left // p.right = a
	if x.left != nil {      // if a != nil
		x.left.parent = x.parent // a.parent = p
	}
	x.parent.parent = x // p.parent = x
	x.left = x.parent   // x.left = p

	x.parent = nil
	st.root = x
}

/*
ggp (x.parent.parent.parent) / gp (x.parent.parent) / p (x.parent)
zigzig: left left case (x is left child's left child)
right rotate (gp) ----->  right rotate (p)
        ggp (1. ggp is nil / 2. gp is left child of ggp / 3. gp is right child of ggp)
       /
      gp                 p                     x
     /  \              /   \                  / \
    p    d            x     gp               a   p
   / \               / \    / \                 / \
  x   c     ----->  a   b  c   d    ----->     b   gp
 / \                                             /  \
a   b                                           c    d
*/
func (st *SplayTree) zigzig(x *node) {
	ggp := x.parent.parent.parent

	isLeftChild := false
	if ggp != nil {
		isLeftChild = ggp.left == x.parent.parent
	}

	x.parent.parent.left = x.parent.right // gp.left = c
	if x.parent.right != nil {            // if c != nil
		x.parent.right.parent = x.parent.parent // c.parent = gp
	}
	x.parent.left = x.right // p.left = b
	if x.right != nil {     // if b != nil
		x.right.parent = x.parent // b.parent = p
	}
	x.parent.right = x.parent.parent  // p.right = gp
	x.parent.parent.parent = x.parent // gp.parent = p
	x.right = x.parent                // x.right = p
	x.parent.parent = x               // p.parent = x
	x.parent = ggp                    // x.parent = ggp

	if ggp == nil {
		st.root = x
	} else if isLeftChild {
		ggp.left = x
	} else {
		ggp.right = x
	}
}

/*
ggp (x.parent.parent.parent) / gp (x.parent.parent) / p (x.parent)
zagzag: right right case (x is right child's right child)
left rotate (gp) ----->  left rotate (p)
        ggp (1. ggp is nil / 2. gp is left child of ggp / 3. gp is right child of ggp)
       /
      gp                      p                      x
     /  \                   /    \                  / \
    d    p                 gp     x                p   b
        / \               /  \   / \              / \
       c   x     ----->  d    c a   b    ----->  gp  a
          / \                                   /  \
         a   b                                 d    c
*/
func (st *SplayTree) zagzap(x *node) {
	ggp := x.parent.parent.parent

	isLeftChild := false
	if ggp != nil {
		isLeftChild = ggp.left == x.parent.parent
	}

	x.parent.parent.right = x.parent.left // gp.right = c
	if x.parent.left != nil {             // if c != nil
		x.parent.left.parent = x.parent.parent // c.parent = gp
	}
	x.parent.right = x.left // p.right = a
	if x.left != nil {      // if a != nil
		x.left.parent = x.parent // a.parent = p
	}
	x.parent.left = x.parent.parent   // p.left = gp
	x.parent.parent.parent = x.parent // gp.parent = p
	x.left = x.parent                 // x.left = p
	x.parent.parent = x               // p.parent = x
	x.parent = ggp                    // x.parent = ggp

	if ggp == nil {
		st.root = x
	} else if isLeftChild {
		ggp.left = x
	} else {
		ggp.right = x
	}
}

/*
ggp (x.parent.parent.parent) / gp (x.parent.parent) / p (x.parent)
zigzag: right left case (x is right child's left child)
right rotate (p) ----->  left rotate (x)
        ggp (1. ggp is nil / 2. gp is left child of ggp / 3. gp is right child of ggp)
       /
      gp                 gp                      x
     /  \               /  \                   /    \
    d    p             d    x                 gp     p
        / \                / \               /  \   / \
       x   c     ----->   a   p    ----->   d    a b   c
      / \                    / \
     a   b                  b   c
*/
func (st *SplayTree) zigzag(x *node) {
	ggp := x.parent.parent.parent

	isLeftChild := false
	if ggp != nil {
		isLeftChild = ggp.left == x.parent.parent
	}

	x.parent.parent.right = x.left // gp.right = a
	if x.left != nil {             // if a != nil
		x.left.parent = x.parent.parent // a.parent = gp
	}
	x.parent.left = x.right // p.left = b
	if x.right != nil {     // if b != nil
		x.right.parent = x.parent // b.parent = p
	}
	x.right = x.parent         // x.right = p
	x.left = x.parent.parent   // x.left = gp
	x.parent.parent.parent = x // gp.parent = x
	x.parent.parent = x        // p.parent = x
	x.parent = ggp             // x.parent = ggp

	if ggp == nil {
		st.root = x
	} else if isLeftChild {
		ggp.left = x
	} else {
		ggp.right = x
	}
}

/*
ggp (x.parent.parent.parent) / gp (x.parent.parent) / p (x.parent)
zagzig: left right case (x is left child's right child)
left rotate (x) ----->  right rotate (g)
        ggp (1. ggp is nil / 2. gp is left child of ggp / 3. gp is right child of ggp)
       /
      gp                 gp                  x
     /  \               /  \               /    \
    p    d             x    d             p     gp
   / \                / \                / \   /  \
  c   x     ----->   p   b    ----->    c   a b    d
     / \            / \
    a   b          c   a
*/
func (st *SplayTree) zagzig(x *node) {
	ggp := x.parent.parent.parent

	isLeftChild := false
	if ggp != nil {
		isLeftChild = ggp.left == x.parent.parent
	}

	x.parent.parent.left = x.right // gp.left = b
	if x.right != nil {            // if b != nil
		x.right.parent = x.parent.parent // b.parent = gp
	}
	x.parent.right = x.left // p.right = a
	if x.left != nil {      // if a != nil
		x.left.parent = x.parent // a.parent = p
	}
	x.left = x.parent          // x.left = p
	x.right = x.parent.parent  // x.right = gp
	x.parent.parent.parent = x // gp.parent = x
	x.parent.parent = x        // p.parent = x
	x.parent = ggp             // x.parent = ggp

	if ggp == nil {
		st.root = x
	} else if isLeftChild {
		ggp.left = x
	} else {
		ggp.right = x
	}
}

func (st *SplayTree) splay(x *node) {
	if x.parent == nil { // x is root
		return
	}
	for x.parent != nil {
		if x.parent.parent == nil && x == x.parent.left { // left
			st.zig(x)
		} else if x.parent.parent == nil && x == x.parent.right { // right
			st.zag(x)
		} else if x == x.parent.left && x.parent == x.parent.parent.left { // left left
			st.zigzig(x)
		} else if x == x.parent.right && x.parent == x.parent.parent.right { // right right
			st.zagzap(x)
		} else if x == x.parent.left && x.parent == x.parent.parent.right { // left right
			st.zigzag(x)
		} else { // right left
			st.zagzig(x)
		}
	}
}

// Search returns true if found key in splay tree
func (st *SplayTree) Search(k interface{}) bool {
	return st.search(st.root, k) != nil
}

func (st *SplayTree) search(x *node, k interface{}) *node {
	if x == nil {
		return nil
	}

	switch st.comparator(k, x.key) {
	case -1:
		return st.search(x.left, k)
	case 0:
		return x
	case 1:
		return st.search(x.right, k)
	default:
		return nil
	}
}

// Get returns value stored inside splay tree
func (st *SplayTree) Get(k interface{}) interface{} {
	return st.search(st.root, k).value
}

// Insert inserts key into splay tree
//	leftMost -> min, rightMost -> max
func (st *SplayTree) Insert(k, v interface{}) error {
	if st.Search(k) {
		return KeyAlreadyExistError{key: k}
	}

	x := st.insert(st.root, k, v)
	st.splay(x)
	st.size++
	return nil
}

func (st *SplayTree) insert(x *node, k, v interface{}) *node {
	if x == nil {
		nr := newNode(k, v)
		st.root = nr
		return nr
	}

	switch st.comparator(k, x.key) {
	case -1:
		if x.left == nil {
			x.left = newNode(k, v)
			x.left.parent = x
			return x.left
		}
		return st.insert(x.left, k, v)
	case 1:
		if x.right == nil {
			x.right = newNode(k, v)
			x.right.parent = x
			return x.right
		}
		return st.insert(x.right, k, v)
	default:
		return nil
	}
}

// Delete deletes k from splay tree
func (st *SplayTree) Delete(k interface{}) error {
	x := st.search(st.root, k)
	if x == nil {
		return KeyNotExistError{key: k}
	}

	st.splay(x)

	if x.left == nil {
		st.replace(x, x.right)
	} else if x.right == nil {
		st.replace(x, x.left)
	} else {
		y := leftMost(x.right)
		if y.parent != x {
			st.replace(y, y.right)
			y.right = x.right
			y.right.parent = y
		}
		st.replace(x, y)
		y.left = x.left
		y.left.parent = y
	}
	st.size--
	return nil
}

func (st *SplayTree) replace(x, y *node) {
	if x.parent == nil {
		st.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	if y != nil {
		y.parent = x.parent
	}
	x = nil
}

/*
// rightMost returns max
func rightMost(x *node) *node {
	if x.right == nil {
		return x
	}
	return rightMost(x.right)
}
*/

// leftMost returns min
func leftMost(x *node) *node {
	if x.left == nil {
		return x
	}
	return leftMost(x.left)
}

// Update updates k, v if k in splay tree
func (st *SplayTree) Update(k, v interface{}) error {
	err := st.Delete(k)
	if err != nil {
		return err
	}

	err = st.Insert(k, v)
	if err != nil {
		return err
	}
	return nil
}

// Empty returns true if splay tree is empty
func (st *SplayTree) Empty() bool {
	return st.size == 0
}

// Size returns quantity of k, v pairs inside splay tree
func (st *SplayTree) Size() int {
	return st.size
}

// Clear clears splay tree
func (st *SplayTree) Clear() {
	st.root = nil
	st.size = 0
}

// Keys returns keys inside splay tree
//	pre-order
func (st *SplayTree) Keys() []interface{} {
	keys := make([]interface{}, 0, st.size)
	st.key(st.root, &keys)
	return keys
}

func (st *SplayTree) key(x *node, keys *[]interface{}) {
	if x == nil {
		return
	}
	*keys = append(*keys, x.key)
	st.key(x.left, keys)
	st.key(x.right, keys)
}

// Values returns keys inside splay tree
//	pre-order
func (st *SplayTree) Values() []interface{} {
	values := make([]interface{}, 0, st.size)
	st.value(st.root, &values)
	return values
}

func (st *SplayTree) value(x *node, values *[]interface{}) {
	if x == nil {
		return
	}
	*values = append(*values, x.value)
	st.value(x.left, values)
	st.value(x.right, values)
}
