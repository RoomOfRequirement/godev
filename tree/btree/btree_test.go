package btree

import (
	"fmt"
	"goContainer"
	"goContainer/tree"
	"goContainer/utils"
	"strconv"
	"testing"
)

func TestBTree(t *testing.T) {
	var _ tree.Tree = (*BTree)(nil)

	bTree := NewBTree(3, container.IntComparator)
	if !bTree.Empty() || bTree.Size() != 0 {
		t.Fail()
	}

	a := []int{1, 2, 3, 4, 5, 6, 7}
	for _, v := range a {
		bTree.Insert(&Item{
			Key:   v,
			Value: fmt.Sprintf("[%s]", strconv.Itoa(v)),
		})
	}
	if bTree.Empty() || bTree.Size() != 7 {
		t.Fail()
	}

	for i, v := range a {
		value, found := bTree.Lookup(v)
		if !found || value.(string) != fmt.Sprintf("[%s]", strconv.Itoa(v)) {
			fmt.Printf("At index %d, value %d\n", i, v)
			t.Fail()
		}
	}

	for i, v := range bTree.Values() {
		if v.(string) != fmt.Sprintf("[%s]", strconv.Itoa(a[i])) {
			fmt.Printf("At index %d, value %s\n", i, v)
			t.Fail()
		}
	}

	for _, v := range a[:3] {
		bTree.Delete(v)
	}

	for i, v := range bTree.Values() {
		if v.(string) != fmt.Sprintf("[%s]", strconv.Itoa(a[3:][i])) {
			fmt.Printf("At index %d, value %s\n", i, v)
			t.Fail()
		}
	}

	if bTree.Size() != 4 {
		t.Fail()
	}

	bTree.Clear()
	if !bTree.Empty() || bTree.Size() != 0 {
		t.Fail()
	}
}

// BenchmarkBTree_Insert-8   	 1000000	      1940 ns/op
func BenchmarkBTree_Insert(b *testing.B) {
	bTree := NewBTree(10, container.IntComparator)
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bTree.Insert(&Item{
			Key:   data[i],
			Value: data[i],
		})
	}
}
