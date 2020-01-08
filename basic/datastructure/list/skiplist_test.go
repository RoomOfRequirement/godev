package list

import (
	"goContainer/basic"
	"goContainer/utils"
	"math"
	"testing"
)

func TestNewSkipList(t *testing.T) {
	var _ basic.Container = (*SkipList)(nil)

	sl := NewSkipList(-10, 0, basic.IntComparator)
	if !sl.Empty() || sl.maxLevel != 0 || sl.decFactor != 2 || sl.Size() != 0 {
		t.Fail()
	}
}

func TestSkipList_Set(t *testing.T) {
	sl := NewSkipList(8, 4, basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		sl.Set(i, a[i])
	}

	if sl.Size() != 8 {
		t.Fail()
	}

	for i, v := range sl.Values() {
		if v.(int) != a[i] {
			t.Fail()
		}
	}
}

func TestSkipList_Delete(t *testing.T) {
	sl := NewSkipList(8, 4, basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		sl.Set(i, a[i])
	}

	v, ok := sl.Delete(1)
	if !ok || v.(int) != 7 {
		t.Fail()
	}

	sl.Set(1, 7)

	for i, n := range a {
		v, ok := sl.Delete(i)
		if !ok || v.(int) != n {
			t.Fail()
		}
	}
}

func TestSkipList_Search(t *testing.T) {
	sl := NewSkipList(8, 4, basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		sl.Set(i, a[i])
	}

	for i := range a {
		if !sl.Search(i) {
			t.Fail()
		}
	}

	if sl.Search(100) {
		t.Fail()
	}
}

func TestSkipList_Get(t *testing.T) {
	sl := NewSkipList(8, 4, basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		sl.Set(i, a[i])
	}

	for i := range a {
		v, found := sl.Get(i)
		if !found || v.(int) != a[i] {
			t.Fail()
		}
	}

	if v, found := sl.Get(100); found || v != nil {
		t.Fail()
	}

	sl.Clear()

	if !sl.Empty() {
		t.Fail()
	}
}

// BenchmarkSkipList_Set-4   	 2000000	       805 ns/op
func BenchmarkSkipList_Set(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	maxLevel := int(math.Log2(float64(b.N)))
	b.ResetTimer()

	sl := NewSkipList(maxLevel, 4, basic.IntComparator)
	for i := 0; i < b.N; i++ {
		sl.Set(i, data[i])
	}
}

// BenchmarkSkipList_Delete-4   	 3000000	       467 ns/op
func BenchmarkSkipList_Delete(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	maxLevel := int(math.Log2(float64(b.N)))

	sl := NewSkipList(maxLevel, 4, basic.IntComparator)
	for i := range data {
		sl.Set(i, data[i])
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(i)
	}
}

// BenchmarkSkipList_Get-4   	 2000000	       650 ns/op
func BenchmarkSkipList_Get(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	maxLevel := int(math.Log2(float64(b.N)))

	sl := NewSkipList(maxLevel, 4, basic.IntComparator)
	for i := range data {
		sl.Set(i, data[i])
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(i)
	}
}
