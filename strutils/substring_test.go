package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	ss := []string{"a", "b", "asd"}
	for _, s := range ss {
		assert.True(t, Contains(s, ss))
	}
	assert.False(t, Contains("hello", ss))
}

func TestLongestCommonSubstring(t *testing.T) {
	assert.Equal(t, "", LongestCommonSubstring("", ""))
	assert.Equal(t, "", LongestCommonSubstring("", "a"))
	assert.Equal(t, "", LongestCommonSubstring("a", ""))
	s1 := "hello"
	s2 := "hello world"
	assert.Equal(t, "hello", LongestCommonSubstring(s1, s2))
}
