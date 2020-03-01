package kmp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	pattern := "ABABC"
	txt := "ABABDABABCABABCABCABABC"
	k := newKMP(pattern)
	idx := k.search(txt)
	assert.Equal(t, 5, idx)

	kmp := New(pattern)
	idx = kmp.SearchFirst(txt)
	assert.Equal(t, 5, idx)
	idx = kmp.SearchLast(txt)
	assert.Equal(t, 18, idx)
	assert.Equal(t, []int{5, 10, 18}, kmp.Search(txt))

	// not exist
	k.search("DDD")
	kmp.Search("DDD")
	kmp.SearchFirst("DDD")
	kmp.SearchLast("DDD")
}
