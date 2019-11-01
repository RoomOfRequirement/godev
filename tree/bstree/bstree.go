package bstree

import (
	"fmt"
	"goContainer"
)

// BSTree Binary Search Tree
//	http://cslibrary.stanford.edu/110/BinaryTrees.html
type BSTree struct {
	Root       *Node
	Comparator container.Comparator
	Diffidence func(a, b interface{}) interface{}
}

// NewBSTree create a nil BSTree
func NewBSTree() *BSTree {
	return &BSTree{
		Root:       nil,
		Comparator: nil,
	}
}

// Node node of BSTree
type Node struct {
	data        interface{}
	left, right *Node
}

// NewNode creates a Node from input data
func NewNode(data interface{}) *Node {
	return &Node{
		data:  data,
		left:  nil,
		right: nil,
	}
}

func (bst *BSTree) lookup(node *Node, target interface{}) bool {
	if node == nil {
		return false
	}

	switch bst.Comparator(target, node.data) {
	case 0:
		return true
	case -1:
		return bst.lookup(node.left, target)
	case 1:
		return bst.lookup(node.right, target)
	default:
		return false
	}
}

// Lookup returns true if target value stored inside the tree
func (bst *BSTree) Lookup(target interface{}) bool {
	return bst.lookup(bst.Root, target)
}

func (bst *BSTree) insert(node *Node, data interface{}) *Node {
	if node == nil {
		return NewNode(data)
	}

	switch bst.Comparator(data, node.data) {
	case -1, 0:
		node.left = bst.insert(node.left, data)
	case 1:
		node.right = bst.insert(node.right, data)
	}
	return node
}

// Insert inserts data into the tree
func (bst *BSTree) Insert(data interface{}) {
	bst.Root = bst.insert(bst.Root, data)
}

func (bst *BSTree) print(node *Node) {
	if node == nil {
		return
	}
	bst.print(node.left)
	fmt.Printf("%+v ", node.data)
	bst.print(node.right)
}

// Print prints the tree with values in increasing order
func (bst *BSTree) Print() {
	bst.print(bst.Root)
	fmt.Println()
}

// Empty returns true if the tree has no value inside (root is nil)
func (bst *BSTree) Empty() bool {
	return bst.Root == nil
}

func (bst *BSTree) size(node *Node) int {
	if node == nil {
		return 0
	}
	return bst.size(node.left) + 1 + bst.size(node.right)
}

// Size returns the number of nodes inside the tree
func (bst *BSTree) Size() int {
	return bst.size(bst.Root)
}

func (bst *BSTree) value(node *Node, dataSlice *[]interface{}) {
	if node == nil {
		return
	}
	bst.value(node.left, dataSlice)
	*dataSlice = append(*dataSlice, node.data)
	bst.value(node.right, dataSlice)
}

// Values returns values stored inside the tree in increasing order
func (bst *BSTree) Values() []interface{} {
	if bst.Root == nil {
		return nil
	}
	var dataSlice []interface{}
	bst.value(bst.Root, &dataSlice)
	return dataSlice
}

// Clear clears the tree by setting root to nil
func (bst *BSTree) Clear() {
	bst.Root = nil
}

func (bst *BSTree) maxDepth(node *Node) int {
	if node == nil {
		return 0
	}
	lDepth := bst.maxDepth(node.left)
	rDepth := bst.maxDepth(node.right)
	if lDepth > rDepth {
		return lDepth + 1
	}
	return rDepth + 1
}

// MaxDepth returns the tree height
func (bst *BSTree) MaxDepth() int {
	return bst.maxDepth(bst.Root)
}

func (bst *BSTree) minValue(node *Node) interface{} {
	currentNode := node
	for currentNode.left != nil {
		currentNode = currentNode.left
	}
	return currentNode.data
}

// MinValue returns the minimum value stored inside the tree
func (bst *BSTree) MinValue() interface{} {
	return bst.minValue(bst.Root)
}

func (bst *BSTree) hasPathSum(node *Node, sum interface{}) bool {
	if node == nil {
		return sum == 0
	}
	sum = bst.Diffidence(sum, node.data)
	return bst.hasPathSum(node.left, sum) || bst.hasPathSum(node.right, sum)
}

// HasPathSum returns true if there's one path to get the input sum by summing up all values along this path
func (bst *BSTree) HasPathSum(sum interface{}) bool {
	if bst.Diffidence == nil {
		panic("no Difference function for node data")
	}
	return bst.hasPathSum(bst.Root, sum)
}

func (bst *BSTree) printPaths(node *Node, path []interface{}) {
	if node == nil {
		return
	}
	path = append(path, node.data)
	if node.left == nil && node.right == nil {
		for _, p := range path {
			fmt.Printf("%+v ", p)
		}
		fmt.Println()
	} else {
		bst.printPaths(node.left, path)
		bst.printPaths(node.right, path)
	}
}

// PrintPaths print all paths from root to leaf
func (bst *BSTree) PrintPaths() {
	var paths []interface{}
	bst.printPaths(bst.Root, paths)
}

func (bst *BSTree) mirror(node *Node) {
	if node == nil {
		return
	}

	bst.mirror(node.left)
	bst.mirror(node.right)
	// swap
	node.left, node.right = node.right, node.left
}

// Mirror mirrors the tree (in-place change)
func (bst *BSTree) Mirror() {
	bst.mirror(bst.Root)
}

func (bst *BSTree) doubleTree(node *Node) {
	if node == nil {
		return
	}

	bst.doubleTree(node.left)
	bst.doubleTree(node.right)

	oldLeft := node.left
	node.left = NewNode(node.data)
	node.left.left = oldLeft
}

// DoubleTree doubles the tree by duplicating nodes and insert them into the tree
func (bst *BSTree) DoubleTree() {
	bst.doubleTree(bst.Root)
}

func (bst *BSTree) sameTree(nodeA, nodeB *Node) bool {
	if nodeA == nil && nodeB == nil {
		return true
	} else if nodeA != nil && nodeB != nil {
		return nodeA.data == nodeB.data && bst.sameTree(nodeA.left, nodeB.left) && bst.sameTree(nodeA.right, nodeB.right)
	}
	return false
}

// SameTree returns true if both trees are made of nodes with the same values
func (bst *BSTree) SameTree(bstB *BSTree) bool {
	return bst.sameTree(bst.Root, bstB.Root)
}

// CountTrees counts how many structurally unique binary search trees which can store the given number of distinct values
func CountTrees(numKeys int) int {
	if numKeys <= 1 {
		return 1
	}
	sum, left, right := 0, 0, 0
	for root := 0; root < numKeys; root++ {
		// left tree node num
		left = CountTrees(root)
		// right tree node num
		right = CountTrees(numKeys - 1 - root)
		sum += left * right
	}
	return sum
}

func (bst *BSTree) isBST(node *Node, min, max interface{}) bool {
	if node == nil {
		return true
	}

	if bst.Comparator(node.data, min) < 0 || bst.Comparator(node.data, max) > 0 {
		return false
	}

	return bst.isBST(node.left, min, node.data) && bst.isBST(node.right, node.data, max)
}

// IsBST returns true if the tree is binary search tree, its inputs should be the min and max values inside the tree
func (bst *BSTree) IsBST(min, max interface{}) bool {
	return bst.isBST(bst.Root, min, max)
}
