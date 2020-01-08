package bag

import (
	"fmt"
	"goContainer/basic"
	"testing"
)

func TestNewBag(t *testing.T) {
	var _ basic.Container = (*Bag)(nil)

	bag := NewBag()

	if !bag.Empty() || len(bag.Entries()) != 0 || len(bag.Values()) != 0 || len(bag.EntriesWithCount()) != 0 {
		t.Fail()
	}

	if bag.Contains(1) || bag.Count(1) != 0 {
		t.Fail()
	}

	if bag.DeleteOne(1) || bag.DeleteAll(1) {
		t.Fail()
	}
}

func TestBag_Add(t *testing.T) {
	bag := NewBag()

	a := []int{1, 2, 2, 3, 3, 3}
	ae := []int{1, 2, 3}

	for i := range a {
		bag.Add(a[i])
	}

	if bag.Size() != 6 {
		t.Fail()
	}

	for i := range ae {
		if !bag.Contains(ae[i]) || bag.Count(ae[i]) != ae[i] {
			t.Fail()
		}
	}

	bag.SetCount(1, 10)
	if bag.Count(1) != 10 {
		t.Fail()
	}

	bag.SetCount(10, 2)
	if !bag.Contains(10) || bag.Count(10) != 2 {
		t.Fail()
	}

	bag.Clear()
	if !bag.Empty() || len(bag.Entries()) != 0 || len(bag.Values()) != 0 || len(bag.EntriesWithCount()) != 0 {
		t.Fail()
	}
}

func TestBag_Delete(t *testing.T) {
	bag := NewBag()

	a := []int{1, 2, 2, 3, 3, 3}

	for i := range a {
		bag.Add(a[i])
	}

	if ok := bag.DeleteOne(3); !ok || bag.Count(3) != 2 {
		t.Fail()
	}

	if ok := bag.DeleteOne(1); !ok || bag.Count(1) != 0 || bag.Contains(1) {
		t.Fail()
	}

	if ok := bag.DeleteAll(2); !ok || bag.Contains(2) {
		t.Fail()
	}
}

func TestBag_Entries(t *testing.T) {
	bag := NewBag()

	a := []int{1, 2, 2, 3, 3, 3}

	for i := range a {
		bag.Add(a[i])
	}

	ae := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}

	es := bag.Entries()
	if len(es) != 3 {
		t.Fail()
	}
	for _, v := range es {
		if ae[v.(int)] != v.(int) {
			t.Fail()
		}
	}

	values := bag.Values()
	if len(values) != 3 {
		t.Fail()
	}
	for _, v := range values {
		if ae[v.(int)] != v.(int) {
			t.Fail()
		}
	}

	esc := bag.EntriesWithCount()
	if len(esc) != 3 {
		t.Fail()
	}
	for _, v := range esc {
		if ae[v.GetEntry().(int)] != v.GetCount() {
			t.Fail()
		}
	}

	bag.ForEachEntry(func(i interface{}) {
		fmt.Print(i)
	})
}

func TestBag_ContainsAll(t *testing.T) {
	bag, other, another := NewBag(), NewBag(), NewBag()

	a := []int{1, 2, 2, 3, 3, 3}
	b := []int{1, 2, 1, 2, 1, 2}
	c := []int{1, 2, 5, 1, 2, 5}

	for i := range a {
		bag.Add(a[i])
		other.Add(b[i])
		another.Add(c[i])
	}

	if !bag.ContainsAll(other) || bag.ContainsAll(another) || !another.ContainsAll(other) {
		t.Fail()
	}

	bag.Merge(another)
	if !bag.ContainsAll(other) {
		t.Fail()
	}
	fmt.Println(bag)
}
