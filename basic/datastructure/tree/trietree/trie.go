package trietree

import (
	"strings"
)

// hash map
//	https://en.wikipedia.org/wiki/Suffix_tree
type tNode struct {
	value    interface{}
	children map[string]*tNode
}

func segFunc(key string, separateRune rune, startIdx int) (segmentedString string, nextIdx int) {
	// validation
	if len(key) == 0 || startIdx < 0 || startIdx > len(key)-1 {
		return "", -1
	}
	nextStart := strings.IndexRune(key[startIdx+1:], separateRune)
	// -1 if separateRune is not present in key -> key string is not able to segment
	if nextStart == -1 {
		return key[startIdx:], -1
	}
	return key[startIdx : startIdx+1+nextStart], startIdx + 1 + nextStart
}

type trie struct {
	segFunc SegmentFunc
	segRune rune
	root    *tNode

	cnt int
}

func newTrie(segRune rune) *trie {
	return &trie{
		segFunc: segFunc,
		segRune: segRune,
		root: &tNode{
			value:    nil,
			children: make(map[string]*tNode),
		},
		cnt: 0, // always has one root node, but not count it
	}
}

func (t *trie) Get(key string) interface{} {
	node := t.root
	r := t.segRune
	// recursively
	for sub, idx := t.segFunc(key, r, 0); sub != ""; sub, idx = t.segFunc(key, r, idx) {
		node = node.children[sub]
		// not present
		if node == nil {
			return nil
		}
	}
	return node.value
}

func (t *trie) Put(key string, value interface{}) (newlyCreated bool) {
	node := t.root
	r := t.segRune
	// recursively
	for sub, idx := t.segFunc(key, r, 0); sub != ""; sub, idx = t.segFunc(key, r, idx) {
		child := node.children[sub]
		// create new path
		if child == nil {
			if node.children == nil {
				node.children = make(map[string]*tNode)
			}
			child = &tNode{}
			node.children[sub] = child

			t.cnt++
		}
		node = child
	}
	if node.value == nil {
		newlyCreated = true
	}
	node.value = value
	return
}

// Delete deletes node's value with key
//	if the node is leaf, delete the node as well
//	need to check the whole search path
func (t *trie) Delete(key string) (found bool) {
	node := t.root
	r := t.segRune
	// record search path
	var path []struct {
		key  string
		node *tNode
	}
	// recursively
	for sub, idx := t.segFunc(key, r, 0); sub != ""; sub, idx = t.segFunc(key, r, idx) {
		path = append(path, struct {
			key  string
			node *tNode
		}{key: sub, node: node})
		node = node.children[sub]
		// not present
		if node == nil {
			return false
		}
	}
	node.value = nil
	t.cnt--
	// leaf node requires to clear the search path
	if isLeaf(node) {
		// reverse iterate
		for i := len(path) - 1; i > -1; i-- {
			pNode := path[i].node
			k := path[i].key
			delete(pNode.children, k)
			t.cnt--

			if isLeaf(pNode) {
				pNode.children = nil
				// pNode value not nil, need to keep
				if pNode.value != nil {
					break
				}
			} else {
				// non-leaf, need to keep
				break
			}
		}
	}
	return true
}

func (t *trie) Walk(wf WalkFunc) error {
	return t.root.walk("", wf)
}

func (t *trie) Size() int {
	return t.cnt
}

func (t *trie) Empty() bool {
	return isLeaf(t.root)
}

func (t *trie) Clear() {
	t.root = &tNode{
		value:    nil,
		children: make(map[string]*tNode),
	}
}

func (t *trie) Values() []interface{} {
	var values []interface{}
	wf := func(values *[]interface{}) WalkFunc {
		return func(key string, value interface{}) error {
			*values = append(*values, value)
			return nil
		}
	}
	_ = t.Walk(wf(&values))
	return values
}

func (n *tNode) walk(key string, wf WalkFunc) error {
	// skip internal created path
	if n.value != nil {
		if err := wf(key, n.value); err != nil {
			return err
		}
	}
	for sub, child := range n.children {
		if err := child.walk(key+sub, wf); err != nil {
			return err
		}
	}
	return nil
}

func isLeaf(node *tNode) bool {
	return len(node.children) == 0
}
