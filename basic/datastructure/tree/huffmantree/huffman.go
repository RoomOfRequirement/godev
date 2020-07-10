package huffmantree

import (
	"fmt"
	"sort"
	"strconv"
)

// Node of huffman tree
type Node struct {
	Left, Right, Parent *Node
	Weight              int
	Value               interface{}
}

// HuffmanTree built on leaves and returns root node
// reference: https://www.wikiwand.com/en/Huffman_coding
type HuffmanTree struct {
	root *Node
	// to hold nodes ptr array
	nodes []*Node
	// for sorting and construct tree
	buf []*Node
}

// NewHuffmanTree ...
func NewHuffmanTree() *HuffmanTree {
	return &HuffmanTree{
		root:  nil,
		nodes: []*Node{},
		buf: []*Node{},
	}
}

// AddNode adds node to huffman tree
//	Notice: call hf.Build() every time you add new node to the tree
func (hf *HuffmanTree) AddNode(node *Node) {
	hf.nodes = append(hf.nodes, node)
}

// AddNodes adds nodes to huffman tree
//	Notice: call hf.Build() every time you add new node to the tree
func (hf *HuffmanTree) AddNodes(nodes ...*Node) {
	for i := range nodes {
		hf.nodes = append(hf.nodes, nodes[i])
	}
}

// Build builds huffman tree
func (hf *HuffmanTree) Build() {
	// make a copy, since i don't want to change nodes array in-place
	hf.buf = hf.nodes[:]
	// empty
	if len(hf.buf) == 0 {
		return
	}
	// sort first
	sort.Stable(SortedNodes(hf.buf))
	// construct tree
	for (len(hf.buf)) > 1 {
		l, r := hf.buf[0], hf.buf[1]
		p := &Node{
			Left:   l,
			Right:  r,
			Parent: nil,
			Weight: l.Weight + r.Weight,
			Value:  nil,
		}
		l.Parent = p
		r.Parent = p
		remained := hf.buf[2:]
		// find parent idx to remained
		idx := sort.Search(len(remained), func(i int) bool {
			return remained[i].Weight >= p.Weight
		})
		// idx + 2 -> idx in hf.nodes
		idx += 2
		// insert parent
		copy(hf.buf[1:], hf.buf[2:idx])
		hf.buf[idx-1] = p
		hf.buf = hf.buf[1:]
	}
	hf.root = hf.buf[0]
}

// PrintTree prints the whole tree with every leaf node's value and code
func (hf *HuffmanTree) PrintTree() {
	if hf.root == nil {
		return
	}
	var traverse func(n *Node, code uint64, bits uint)
	traverse = func(n *Node, code uint64, bits uint) {
		// Print only leaf
		if n.Left == nil {
			// bits to hold leading 0
			fmt.Printf("'%v': %0"+strconv.Itoa(int(bits))+"b\n", n.Value, code)
			return
		}
		// bit length ++
		bits++
		// left bit 0
		traverse(n.Left, code<<1, bits)
		// right bit 1
		traverse(n.Right, (code<<1)+1, bits)
	}

	traverse(hf.root, 0, 0)
}

// Clear only clears hf's nodes ptr array, not clear every node
func (hf *HuffmanTree) Clear() {
	hf.root = nil
	hf.nodes = hf.nodes[:0]
	hf.buf = hf.buf[:0]
}

// Code returns huffman code of tree node
//	Notice: call it after building huffman tree (hf.Build())
func (n *Node) Code() (code uint64, bits uint) {
	// bottom-up, from leaf node to root
	for parent := n.Parent; parent != nil; n, parent = parent, parent.Parent {
		// node on parent's right leaf -> bit 1
		// left -> bit 0
		if parent.Right == n {
			code |= 1 << bits
		}
		bits++
	}
	return
}

// SortedNodes for sorting
type SortedNodes []*Node

// Len ...
func (sn SortedNodes) Len() int {
	return len(sn)
}

// Less ...
func (sn SortedNodes) Less(i, j int) bool {
	return sn[i].Weight < sn[j].Weight
}

// Swap ...
func (sn SortedNodes) Swap(i, j int) {
	sn[i], sn[j] = sn[j], sn[i]
}
