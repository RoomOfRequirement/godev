package treemap

import (
	"fmt"
	"godev/basic"
	"godev/basic/datastructure/maps"
	"godev/utils"
	"testing"
)

func TestNewMap(t *testing.T) {
	var _ basic.Container = (*Map)(nil)
	var _ maps.Map = (*Map)(nil)

	m := NewMap(basic.IntComparator)

	if !m.Empty() {
		t.Fail()
	}
}

func TestMap_Set(t *testing.T) {
	m := NewMap(basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		m.Set(i, a[i])
	}

	if m.Size() != 8 {
		t.Fail()
	}

	for i, v := range m.Keys() {
		if v.(int) != i {
			t.Fail()
		}
	}

	for i, v := range m.Values() {
		if v.(int) != a[i] {
			t.Fail()
		}
	}

	m.Set(100, 100)
	if v, found := m.Get(100); !found || v.(int) != 100 {
		t.Fail()
	}

	m.Clear()
	if !m.Empty() {
		t.Fail()
	}
}

func TestMap_Delete(t *testing.T) {
	m := NewMap(basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		m.Set(i, a[i])
	}

	for i := range a {
		ok := m.Delete(i)
		if !ok {
			t.Fail()
		}
	}

	if !m.Empty() {
		fmt.Println(m.Size())
		t.Fail()
	}
}

func TestMap_Iterator(t *testing.T) {
	m := NewMap(basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		m.Set(i, a[i])
	}

	it := m.Iterator()
	i := 0
	for it.HasNext() {
		k, v := it.Next()
		if k.(int) != i || v.(int) != a[i] {
			t.Fail()
		}
		i++
	}
}

// BenchmarkMap_Set-8   	 1000000	      1717 ns/op
func BenchmarkMap_Set(b *testing.B) {
	m := NewMap(basic.IntComparator)
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(data[i], data[i])
	}
}

// BenchmarkMap_Get-8   	 1000000	      1197 ns/op
func BenchmarkMap_Get(b *testing.B) {
	m := NewMap(basic.IntComparator)
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
		m.Set(data[i], data[i])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Get(data[i])
	}
}
