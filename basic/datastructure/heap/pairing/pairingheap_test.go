package pairing

import (
	"goContainer/basic/datastructure/heap"
	"goContainer/utils"
	"testing"
)

type item int

func (it item) Compare(ait heap.Item) int {
	if it > ait.(item) {
		return 1
	} else if it == ait.(item) {
		return 0
	} else {
		return -1
	}
}

func TestNewHeap(t *testing.T) {
	var _ heap.Heap = (*Heap)(nil)
	h := NewHeap()
	if h.root != nil || h.itemNum != 0 || !h.Empty() || h.FindMin() != nil || h.DeleteMin() != nil {
		t.Fail()
	}
}

func TestHeap_Insert(t *testing.T) {
	a := []item{12, 7, 25, 15, 28, 33, 41, 1}

	h := NewHeap()

	for i := range a {
		h.Insert(a[i])
	}

	if h.Empty() {
		t.Fail()
	}

	if h.FindMin().(item) != 1 || h.Size() != 8 {
		t.Fail()
	}

	min := h.DeleteMin()
	if min.(item) != 1 || h.FindMin().(item) != 7 || h.Size() != 7 {
		t.Fail()
	}

	h.Clear()

	if !h.Empty() || h.FindMin() != nil || h.DeleteMin() != nil || h.Size() != 0 {
		t.Fail()
	}
}

func TestHeap_Delete(t *testing.T) {
	a := []item{12, 7, 25, 15, 28, 33, 41, 1}

	h := NewHeap()

	for i := range a {
		h.Insert(a[i])
	}

	for i := range a {
		item := h.Delete(a[i])
		if item != a[i] {
			t.Fail()
		}
	}

	if !h.Empty() || h.Size() != 0 {
		t.Fail()
	}
}

func TestHeap_Update(t *testing.T) {
	a := []item{12, 7, 25, 15, 28, 33, 41, 1}

	h := NewHeap()

	for i := range a {
		h.Insert(a[i])
	}

	updated := h.Update(a[1], item(0))
	if !updated || h.FindMin() != item(0) {
		t.Fail()
	}

	updated = h.Update(a[4], item(50))
	if !updated || !h.Search(item(50)) {
		t.Fail()
	}

	updated = h.Update(item(100), item(50))
	if updated {
		t.Fail()
	}
}

func TestHeap_Values(t *testing.T) {
	a := []item{12, 7, 25, 15, 28, 33, 41, 1}

	h := NewHeap()

	expectedMap := make(map[item]struct{}, len(a))

	for i := range a {
		h.Insert(a[i])
		expectedMap[a[i]] = struct{}{}
	}

	for _, v := range h.Values() {
		if _, found := expectedMap[v.(item)]; !found {
			t.Fatal(v)
		}
	}
}

func TestHeap_Meld(t *testing.T) {
	a := []item{12, 7, 25, 15, 28, 33, 41, 1}
	b := []item{18, 35, 20, 42, 9, 31, 23, 6, 48, 11, 24, 52, 13, 2}

	h1, h2 := NewHeap(), NewHeap()
	expectedMap := make(map[item]struct{}, len(a)+len(b))

	for i := range a {
		h1.Insert(a[i])
		expectedMap[a[i]] = struct{}{}
	}
	for i := range b {
		h2.Insert(b[i])
		expectedMap[b[i]] = struct{}{}
	}

	h := h1.Meld(h2)

	if h.Size() != len(expectedMap) {
		t.Fail()
	}

	for _, v := range h.Values() {
		if _, found := expectedMap[v.(item)]; !found {
			t.Fatal(v)
		}
	}
}

// TODO: extremely slow...
func BenchmarkHeap_Insert(b *testing.B) {
	data := make([]item, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = item(utils.GenerateRandomInt())
	}
	b.ResetTimer()

	h := NewHeap()
	for i := 0; i < b.N; i++ {
		h.Insert(data[i])
	}
}
