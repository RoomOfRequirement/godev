package utils

import (
	"math"
	"strconv"
	"testing"
)

func BenchmarkGenerateRandomInt(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				for j := 0; j < n; j++ {
					_ = GenerateRandomInt()
				}
			}
		})
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				_ = GenerateRandomString(n)
			}
		})
	}
}

func TestGenerateRandomInt(t *testing.T) {
	var n interface{}
	n = GenerateRandomInt()
	_, ok := n.(int)
	if !ok {
		t.Fail()
	}
}

func TestGenerateRandomIntInRange(t *testing.T) {
	s := GenerateRandomInt()
	e := s + 1000
	n := GenerateRandomIntInRange(s, e)
	if n < s || n >= e {
		t.Fail()
	}

	n = GenerateRandomIntInRange(e, s)
	if n < s || n >= e {
		t.Fail()
	}
}

func TestGenerateRandomString(t *testing.T) {
	var s interface{}
	s = GenerateRandomString(10)
	str, ok := s.(string)
	if !ok || len(str) != 10 {
		t.Fail()
	}
}
