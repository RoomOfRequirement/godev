package huffmantree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHuffmanTree(t *testing.T) {
	hf := NewHuffmanTree()
	hf.Build()
	assert.Nil(t, hf.root)
	hf.PrintTree()
	hf.AddNode(&Node{Value: "x", Weight: 5})
	hf.Build()
	assert.Equal(t, "x", hf.root.Value)
	assert.Equal(t, 5, hf.root.Weight)
	hf.Clear()
	assert.Nil(t, hf.root)
	leaves := []*Node{
		{Value: "a", Weight: 11},
		{Value: "c", Weight: 2},
		{Value: "b", Weight: 6},
		{Value: "e", Weight: 7},
		{Value: " ", Weight: 10},
		{Value: "d", Weight: 10},
	}
	hf.AddNodes(leaves...)
	hf.Build()
	hf.PrintTree()
	expected := map[string]struct {
		uint64
		uint
	}{
		" ": {0, 2},  // 00
		"d": {1, 2},  // 01
		"a": {2, 2},  // 10
		"e": {6, 3},  // 110
		"c": {14, 4}, // 1110
		"b": {15, 4}, // 1111
	}
	for _, leaf := range leaves {
		code, bits := leaf.Code()
		assert.Equal(t, expected[leaf.Value.(string)].uint64, code)
		assert.Equal(t, expected[leaf.Value.(string)].uint, bits)
	}
	hf.AddNode(&Node{Value: "x", Weight: 5})
	hf.Build()
	hf.PrintTree()
}
