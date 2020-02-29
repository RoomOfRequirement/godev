package quickselect

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

func TestQuickSelect(t *testing.T) {
	arrInt := []int{5, 9, 77, 62, 71, 11, 22, 46, 36, 18, 19, 33, 75, 17, 39, 41, 73, 50, 217, 79, 120}
	is := sort.IntSlice{5, 9, 77, 62, 71, 11, 22, 46, 36, 18, 19, 33, 75, 17, 39, 41, 73, 50, 217, 79, 120}
	isCopy := is[:]
	sort.Sort(isCopy)
	n := is.Len()

	for i := 0; i < n; i++ {
		assert.Equal(t, isCopy[i], is[QuickSelect(is, i)])
		assert.Equal(t, isCopy[i], MOM(arrInt, i))
	}

	k := rand.Intn(n)
	assert.Equal(t, isCopy[k], is[QuickSelect(is, k)])
	assert.Equal(t, isCopy[k], MOM(arrInt, k))
}

func BenchmarkQuickSelect(b *testing.B) {
	is := sort.IntSlice{5, 9, 77, 62, 71, 11, 22, 46, 36, 18, 19, 33, 75, 17, 39, 41, 73, 50, 217, 79, 120}
	arrInt := []int{5, 9, 77, 62, 71, 11, 22, 46, 36, 18, 19, 33, 75, 17, 39, 41, 73, 50, 217, 79, 120}
	k := rand.Intn(is.Len())

	b.Run("QuickSelect", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = QuickSelect(is, k)
		}
	})
	b.Run("MOM", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = MOM(arrInt, k)
		}
	})
}
