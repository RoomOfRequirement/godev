package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	s := "¢ह€한𐍈"
	expected := "𐍈한€ह¢"
	assert.Equal(t, expected, Reverse(s))

	s = "bròwn"
	expected = "nwòrb"
	assert.Equal(t, expected, ReversePreservingCombiningCharacters(s))

	assert.Equal(t, "", ReversePreservingCombiningCharacters(""))
}

func BenchmarkReverse(b *testing.B) {
	s := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		_ = Reverse(s)
	}
}

func BenchmarkReversePreservingCombiningCharacters(b *testing.B) {
	s := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		_ = ReversePreservingCombiningCharacters(s)
	}
}
