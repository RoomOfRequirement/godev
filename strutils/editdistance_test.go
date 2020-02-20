package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEditDistance(t *testing.T) {
	s1 := "crate"
	s2 := "trace"
	d, err := EditDistance(Levenshtein, s1, s2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d)
	d, err = EditDistance(LongestCommonSubsequence, s1, s2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d)
	d, err = EditDistance(Hamming, s1, s2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d)

	// empty
	d, err = EditDistance(Levenshtein, "", s2)
	assert.NoError(t, err)
	assert.Equal(t, len(s2), d)
	d, err = EditDistance(LongestCommonSubsequence, "", s2)
	assert.NoError(t, err)
	assert.Equal(t, len(s2), d)

	// err
	d, err = EditDistance(10, s1, s2)
	assert.Error(t, err)
	assert.Equal(t, -1, d)
	d, err = EditDistance(Hamming, "s1", "s12")
	assert.Error(t, err)
	assert.Equal(t, -1, d)
}
