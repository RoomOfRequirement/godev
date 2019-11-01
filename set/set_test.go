package set

import (
	"goContainer"
	"testing"
)

func TestNewIntSet(t *testing.T) {
	// containerInterfaceAssertion()
	var _ container.Container = (*IntSet)(nil)
	var _ container.Container = (*FloatSet)(nil)

	s := NewIntSet(-2)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	if s.Size() != 3 {
		t.Fail()
	}
}

func TestNewIntSetFromIntArray(t *testing.T) {
	s := NewIntSetFromIntArray([]int{1, 2, 3})
	if !s.Get(1) || !s.Get(2) || !s.Get(3) {
		t.Fail()
	}
}

func TestIntSet_Empty(t *testing.T) {
	s := NewIntSet(10)
	if !s.Empty() {
		t.Fail()
	}
}

func TestIntSet_Clear(t *testing.T) {
	s := NewIntSetFromIntArray([]int{1, 2, 3})
	s.Clear()
	if !s.Empty() {
		t.Fail()
	}
}

func TestIntSet_Values(t *testing.T) {
	a := []int{1, 2, 3}
	s := NewIntSetFromIntArray(a)
	values := container.GetSortedValues(s, container.IntComparator, false)
	for i, v := range a {
		if v != values[i].(int) {
			t.Fail()
		}
	}
}

func TestIntSetEqual(t *testing.T) {
	if !IntSetEqual(NewIntSetFromIntArray([]int{1, 2, 3}), NewIntSetFromIntArray([]int{1, 2, 3})) {
		t.Fail()
	}
	if IntSetEqual(NewIntSetFromIntArray([]int{1, 2, 3}), NewIntSetFromIntArray([]int{2, 3, 5})) {
		t.Fail()
	}
}

func TestIntSet_Add(t *testing.T) {
	s := NewIntSet(1)
	s.Add(1)
	if !s.Get(1) {
		t.Fail()
	}
}

func TestIntSet_Get(t *testing.T) {
	s := NewIntSet(1)
	s.Add(1)
	if !s.Get(1) {
		t.Fail()
	}
}

func TestIntSet_Delete(t *testing.T) {
	s := NewIntSet(1)
	s.Add(1)
	s.Delete(1)
	if s.Get(1) {
		t.Fail()
	}
}

func TestIntSet_Size(t *testing.T) {
	s := NewIntSet(1)
	s.Add(1)
	if s.Size() != 1 {
		t.Fail()
	}
}

func TestIntSet_Union(t *testing.T) {
	s := NewIntSet(1)
	s.Add(1)
	s1 := NewIntSet(1)
	s1.Add(2)
	s2 := s.Union(s1)
	if !s2.Get(1) || !s2.Get(2) {
		t.Fail()
	}
}

func TestIntSet_Intersection(t *testing.T) {
	s := NewIntSetFromIntArray([]int{1, 2, 3})
	s1 := NewIntSetFromIntArray([]int{2, 3, 5})
	s2 := s.Intersection(s1)
	if !IntSetEqual(s2, NewIntSetFromIntArray([]int{2, 3})) {
		t.Fail()
	}
}

func TestIntSet_Difference(t *testing.T) {
	s := NewIntSetFromIntArray([]int{1, 2, 3})
	s1 := NewIntSetFromIntArray([]int{2, 3, 5})
	s2 := s.Difference(s1)
	if !IntSetEqual(s2, NewIntSetFromIntArray([]int{1, 5})) {
		t.Fail()
	}
}

func TestNewFloatSet(t *testing.T) {
	s := NewFloatSet(-2)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	if s.Size() != 3 {
		t.Fail()
	}
}

func TestNewFloatSetFromFloatArray(t *testing.T) {
	s := NewFloatSetFromFloatArray([]float64{1, 2, 3})
	if !s.Get(1) || !s.Get(2) || !s.Get(3) {
		t.Fail()
	}
}

func TestFloatSet_Empty(t *testing.T) {
	s := NewFloatSet(10)
	if !s.Empty() {
		t.Fail()
	}
}

func TestFloatSet_Clear(t *testing.T) {
	s := NewFloatSetFromFloatArray([]float64{1, 2, 3})
	s.Clear()
	if !s.Empty() {
		t.Fail()
	}
}

func TestFloatSet_Values(t *testing.T) {
	a := []float64{1, 2, 3}
	s := NewFloatSetFromFloatArray(a)
	values := container.GetSortedValues(s, container.Float64Comparator, false)
	for i, v := range a {
		if v != values[i].(float64) {
			t.Fail()
		}
	}
}

func TestFloatSetEqual(t *testing.T) {
	if !FloatSetEqual(NewFloatSetFromFloatArray([]float64{1, 2, 3}), NewFloatSetFromFloatArray([]float64{1, 2, 3})) {
		t.Fail()
	}
	if FloatSetEqual(NewFloatSetFromFloatArray([]float64{1, 2, 3}), NewFloatSetFromFloatArray([]float64{2, 3, 5})) {
		t.Fail()
	}
}

func TestFloatSet_Add(t *testing.T) {
	s := NewFloatSet(1)
	s.Add(1)
	if !s.Get(1) {
		t.Fail()
	}
}

func TestFloatSet_Get(t *testing.T) {
	s := NewFloatSet(1)
	s.Add(1)
	if !s.Get(1) {
		t.Fail()
	}
}

func TestFloatSet_Delete(t *testing.T) {
	s := NewFloatSet(1)
	s.Add(1)
	s.Delete(1)
	if s.Get(1) {
		t.Fail()
	}
}

func TestFloatSet_Size(t *testing.T) {
	s := NewFloatSet(1)
	s.Add(1)
	if s.Size() != 1 {
		t.Fail()
	}
}

func TestFloatSet_Union(t *testing.T) {
	s := NewFloatSet(1)
	s.Add(1)
	s1 := NewFloatSet(1)
	s1.Add(2)
	s2 := s.Union(s1)
	if !s2.Get(1) || !s2.Get(2) {
		t.Fail()
	}
}

func TestFloatSet_Intersection(t *testing.T) {
	s := NewFloatSetFromFloatArray([]float64{1, 2, 3})
	s1 := NewFloatSetFromFloatArray([]float64{2, 3, 5})
	s2 := s.Intersection(s1)
	if !FloatSetEqual(s2, NewFloatSetFromFloatArray([]float64{2, 3})) {
		t.Fail()
	}
}

func TestFloatSet_Difference(t *testing.T) {
	s := NewFloatSetFromFloatArray([]float64{1, 2, 3})
	s1 := NewFloatSetFromFloatArray([]float64{2, 3, 5})
	s2 := s.Difference(s1)
	if !FloatSetEqual(s2, NewFloatSetFromFloatArray([]float64{1, 5})) {
		t.Fail()
	}
}
