package bidirectionalmap

import (
	"fmt"
	"goContainer/basic"
	"goContainer/basic/datastructure/maps"
	"goContainer/utils"
	"testing"
)

func TestNewBidirectionalMap(t *testing.T) {
	var _ basic.Container = (*BidirectionalMap)(nil)
	var _ maps.Map = (*BidirectionalMap)(nil)

	m := NewBidirectionalMap()

	if !m.Empty() || m.Size() != 0 || len(m.Keys()) != 0 || len(m.Values()) != 0 {
		t.Fail()
	}
}

func TestBidirectionalMap_Set(t *testing.T) {
	m := NewBidirectionalMap()

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		m.Set(i, a[i])
	}

	if m.Size() != 8 || len(m.Keys()) != 8 || len(m.Values()) != 8 {
		t.Fail()
	}

	for i := range a {
		if v, found := m.Get(i); !found || v.(int) != a[i] {
			t.Fail()
		}
		if v, found := m.Get(a[i]); !found || v.(int) != i {
			t.Fail()
		}
	}

	expectedKeys := make(map[int]struct{}, len(a))
	expectedValues := make(map[int]struct{}, len(a))

	for k, v := range a {
		expectedKeys[k] = struct{}{}
		expectedValues[v] = struct{}{}
	}

	for _, k := range m.Keys() {
		if _, found := expectedKeys[k.(int)]; !found {
			t.Fail()
		}
	}

	for _, k := range m.Values() {
		if _, found := expectedValues[k.(int)]; !found {
			t.Fail()
		}
	}

	m.Set(1, 100)
	if v, found := m.Get(1); !found || v.(int) != 100 {
		t.Fail()
	}

	m.Set(100, 100)
	if v, found := m.Get(100); !found || v.(int) != 100 {
		t.Fail()
	}

	if _, found := m.Get(200); found {
		t.Fail()
	}

	m.Clear()
	if !m.Empty() {
		t.Fail()
	}
}

func TestBidirectionalMap_Delete(t *testing.T) {
	m := NewBidirectionalMap()

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

	if ok := m.Delete(100); ok {
		t.Fail()
	}

	m.Set(2, 25)

	if ok := m.Delete(25); !ok {
		t.Fail()
	}

	if !m.Empty() {
		fmt.Println(m.Size())
		t.Fail()
	}
}

func TestBidirectionalMap_Iterator(t *testing.T) {
	m := NewBidirectionalMap()

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	expected := make(map[int]int, len(a))
	for i := range a {
		m.Set(i, a[i])
		expected[i] = a[i]
	}

	it := m.Iterator()
	i := 0
	for it.HasNext() {
		k, v := it.Next()
		if v.(int) != expected[k.(int)] {
			t.Fail()
		}
		i++
	}
}

// BenchmarkBidirectionalMap_Set-8   	 1000000	      1123 ns/op
func BenchmarkBidirectionalMap_Set(b *testing.B) {
	m := NewBidirectionalMap()
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(data[i], data[i])
	}
}

// BenchmarkBidirectionalMap_Get-8   	10000000	       148 ns/op
func BenchmarkBidirectionalMap_Get(b *testing.B) {
	m := NewBidirectionalMap()
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

// BenchmarkBidirectionalMap_Delete-8   	 5000000	       400 ns/op
func BenchmarkBidirectionalMap_Delete(b *testing.B) {
	m := NewBidirectionalMap()
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
		m.Set(data[i], data[i])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Delete(data[i])
	}
}
