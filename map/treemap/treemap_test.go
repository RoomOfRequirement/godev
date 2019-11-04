package treemap

import (
	"fmt"
	container "goContainer"
	"testing"
)

func TestNewMap(t *testing.T) {
	var _ container.Container = (*Map)(nil)

	m := NewMap(container.IntComparator)

	if !m.Empty() {
		t.Fail()
	}
}

func TestMap_Set(t *testing.T) {
	m := NewMap(container.IntComparator)

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
	m := NewMap(container.IntComparator)

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
	m := NewMap(container.IntComparator)

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
