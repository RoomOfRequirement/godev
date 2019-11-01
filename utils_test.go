package container

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
					GenerateRandomInt()
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
				GenerateRandomString(n)
			}
		})
	}
}
