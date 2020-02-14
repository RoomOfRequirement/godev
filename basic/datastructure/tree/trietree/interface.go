package trietree

import "godev/basic"

// Trie tree
type Trie interface {
	basic.Container

	Get(key string) interface{}
	Put(key string, value interface{}) (newlyCreated bool)
	Delete(key string) (found bool)
	Walk(wf WalkFunc) error
}

// WalkFunc ...
type WalkFunc func(key string, value interface{}) error

// SegmentFunc segments key from startIdx with separateRune
//	it returns segmentedString and nextIdx for next segmentation
type SegmentFunc func(key string, separateRune rune, startIdx int) (segmentedString string, nextIdx int)

// New creates a new Trie
func New(segRune rune) Trie {
	return newTrie(segRune)
}
