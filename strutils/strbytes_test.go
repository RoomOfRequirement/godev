package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafeStringBytesCast(t *testing.T) {
	strs := []string{
		"ashhr5urtrturhYGFJHWGKJDHUHJHNWJ",
		"jopijoper124OUIPIPOKMKNNBNCZZSA",
		"`1234567890-=~!@#$%^&*()_+",
		"{}\t\b[]':;,.<>?/\n",
		"¢ह€한𐍈",
	}
	n := len(strs)
	for i := 0; i < n; i++ {
		s := StringToBytes(strs[i])
		ss := StringToBytes(strs[i])
		ns, nss := BytesToString(s), BytesToString(ss)
		if ns != strs[i] || nss != strs[i] || cap(s) != cap(ss) {
			t.Fail()
		}
	}
	test := struct {
		s1, s2 string
	}{"hello", "world"}
	assert.Equal(t, "hello", BytesToString(StringToBytes(test.s1)))
	assert.Equal(t, "world", BytesToString(StringToBytes(test.s2)))
}

func BenchmarkOfficialCast(b *testing.B) {
	strs := []string{
		"ashhr5urtrturhYGFJHWGKJDHUHJHNWJ",
		"jopijoper124OUIPIPOKMKNNBNCZZSA",
		"`1234567890-=~!@#$%^&*()_+",
		"{}\t\b[]':;,.<>?/\n",
		"¢ह€한𐍈",
	}
	n := len(strs)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			b := []byte(strs[j])
			_ = string(b)
		}
	}
}

func BenchmarkUnsafeCast(b *testing.B) {
	strs := []string{
		"ashhr5urtrturhYGFJHWGKJDHUHJHNWJ",
		"jopijoper124OUIPIPOKMKNNBNCZZSA",
		"`1234567890-=~!@#$%^&*()_+",
		"{}\t\b[]':;,.<>?/\n",
		"¢ह€한𐍈",
	}
	n := len(strs)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			s := StringToBytes(strs[j])
			_ = BytesToString(s)
		}
	}
}
