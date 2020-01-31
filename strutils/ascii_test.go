package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode"
)

func TestIsAllAscii(t *testing.T) {
	assert.True(t, IsAllASCII(ascii()))
	assert.True(t, isAllASCII(ascii()))

	assert.False(t, IsAllASCII("¢ह€한𐍈"))
	assert.False(t, isAllASCII("¢ह€한𐍈"))
}

func BenchmarkIsAllAscii(b *testing.B) {
	str := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := IsAllASCII(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

func BenchmarkIsAscii(b *testing.B) {
	str := ascii()
	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		is := isAllASCII(str)
		if !is {
			b.Log("notASCII")
		}
	}
}

func ascii() string {
	byt := make([]byte, unicode.MaxASCII+1)
	for i := range byt {
		byt[i] = byte(i)
	}
	return string(byt)
}
