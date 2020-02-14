package trietree

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"godev/basic/datastructure/tree"
	"testing"
)

func TestNew(t *testing.T) {
	var _ tree.Tree = (Trie)(nil)
	trie := New('/')

	assert.True(t, trie.Empty())

	newlyCreated := trie.Put("/a", "a")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/a/b", "b")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/a/c", "c")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/b/d", "d")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/c/e", "e")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/c/e", "f")
	assert.False(t, newlyCreated)

	assert.Equal(t, "a", trie.Get("/a"))
	assert.Equal(t, "b", trie.Get("/a/b"))
	assert.Equal(t, "c", trie.Get("/a/c"))
	assert.Equal(t, "d", trie.Get("/b/d"))
	assert.Equal(t, "f", trie.Get("/c/e"))

	// Put "/a/c" as internal node
	trie.Put("/a/c", nil)
	assert.Nil(t, trie.Get("/a/c"))
	assert.Nil(t, trie.Get("/a/c/e"))

	found := trie.Delete("/c/e")
	assert.True(t, found)
	assert.Nil(t, trie.Get("/c/e"))
	assert.Nil(t, trie.Get("/a/c/e"))
	found = trie.Delete("/c/e")
	assert.False(t, found)

	values := trie.Values()
	assert.NotNil(t, values)

	size := trie.Size()
	assert.Equal(t, 4, size)

	assert.False(t, trie.Empty())

	err := trie.Walk(func(key string, value interface{}) error {
		return errors.New("test")
	})
	assert.Error(t, err)

	trie.Clear()
	assert.True(t, trie.Empty())
	assert.Nil(t, trie.Values())

	newlyCreated = trie.Put("/a/b/c/d/e", "x")
	assert.True(t, newlyCreated)
	newlyCreated = trie.Put("/a/b/c/d/e/f", "y")
	assert.True(t, newlyCreated)
	found = trie.Delete("/a/b/c/d/e/f")
	assert.True(t, found)
}
