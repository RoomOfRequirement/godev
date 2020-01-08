package linkedhashmap

import (
	"fmt"
	"goContainer/basic"
	"goContainer/basic/datastructure/maps"
	"goContainer/utils"
	"testing"
)

func TestNewLinkedHashMap(t *testing.T) {
	var _ basic.Container = (*LinkedHashMap)(nil)
	var _ maps.Map = (*LinkedHashMap)(nil)

	m := NewLinkedHashMap()

	if !m.Empty() || m.Size() != 0 || len(m.Keys()) != 0 || len(m.Values()) != 0 {
		t.Fail()
	}
}

func TestLinkedHashMap_Set(t *testing.T) {
	m := NewLinkedHashMap()

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
	}

	for i := range a {
		if m.Keys()[i] != i || m.Values()[i] != a[i] {
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

	m.Clear()
	if !m.Empty() {
		t.Fail()
	}
}

func TestLinkedHashMap_Delete(t *testing.T) {
	m := NewLinkedHashMap()

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

	if !m.Empty() {
		fmt.Println(m.Size())
		t.Fail()
	}
}

func TestLinkedHashMap_Iterator(t *testing.T) {
	m := NewLinkedHashMap()

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

// BenchmarkLinkedHashMap_Set-8   	 2000000	       833 ns/op
func BenchmarkLinkedHashMap_Set(b *testing.B) {
	m := NewLinkedHashMap()
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(data[i], data[i])
	}
}

// BenchmarkLinkedHashMap_Get-8   	10000000	       146 ns/op
func BenchmarkLinkedHashMap_Get(b *testing.B) {
	m := NewLinkedHashMap()
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

// BenchmarkLinkedHashMap_Delete-8   	10000000	       212 ns/op
func BenchmarkLinkedHashMap_Delete(b *testing.B) {
	m := NewLinkedHashMap()
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

func BenchmarkLinkedHashMap(b *testing.B) {
	b.Run("HashMap-Set: ", func(b *testing.B) {
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
		}

		m := make(map[interface{}]interface{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m[data[i]] = data[i]
		}
	})

	b.Run("HashMap-Get: ", func(b *testing.B) {
		data := make([]int, b.N)
		m := make(map[interface{}]interface{})
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
			m[data[i]] = data[i]
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m[data[i]]
		}
	})

	b.Run("HashMap-Delete: ", func(b *testing.B) {
		data := make([]int, b.N)
		m := make(map[interface{}]interface{})
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
			m[data[i]] = data[i]
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			delete(m, data[i])
		}
	})

	b.Run("LinkedHashMap-Set: ", func(b *testing.B) {
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
		}

		m := NewLinkedHashMap()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Set(data[i], data[i])
		}
	})

	b.Run("LinkedHashMap-Get: ", func(b *testing.B) {
		m := NewLinkedHashMap()
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
			m.Set(data[i], data[i])
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m.Get(data[i])
		}
	})

	b.Run("LinkedHashMap-Delete: ", func(b *testing.B) {
		m := NewLinkedHashMap()
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
			m.Set(data[i], data[i])
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m.Delete(data[i])
		}
	})
}

/*
BenchmarkLinkedHashMap/HashMap-Set:_-8         	             2000000	       604 ns/op
BenchmarkLinkedHashMap/LinkedHashMap-Set:_-8   	             2000000	       785 ns/op

BenchmarkLinkedHashMap/HashMap-Get:_-8         	            10000000	       143 ns/op
BenchmarkLinkedHashMap/LinkedHashMap-Get:_-8   	            20000000	       154 ns/op

BenchmarkLinkedHashMap/HashMap-Delete:_-8      	            10000000	       182 ns/op
BenchmarkLinkedHashMap/LinkedHashMap-Delete:_-8         	10000000	       243 ns/op
*/
