package rbtree

import (
	"fmt"
	"goContainer"
	"goContainer/tree"
	"testing"
)

func TestRBTree(t *testing.T) {
	root := NewNode(0)
	rbTree := NewRBTree(root, container.IntComparator)

	var _ tree.Tree = (*RBTree)(nil)

	if rbTree.Empty() {
		t.Fail()
	}

	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, i := range a {
		rbTree.Insert(i)
	}

	if rbTree.Root.value.(int) != 2 {
		t.Fail()
	}
	if rbTree.MinValue().(int) != -8 {
		t.Fail()
	}

	fmt.Println(rbTree)

	ok := rbTree.Delete(0)
	if !ok {
		t.Fail()
	}

	aSorted := []int{-8, -5, -3, 1, 2, 3, 4, 5, 6, 8}
	for i, v := range rbTree.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
	if rbTree.Size() != len(aSorted) {
		t.Fail()
	}

	b := []int{-3, 1, 2, 3}
	for _, v := range b {
		if ok := rbTree.Delete(v); !ok {
			t.Fail()
		}
	}

	aSorted = []int{-8, -5, 4, 5, 6, 8}
	for i, v := range rbTree.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}

	rbTree.Clear()
	if !rbTree.Empty() || rbTree.Size() != 0 || rbTree.Values() != nil {
		t.Fail()
	}
}

// BenchmarkRBTree_Insert-8   	 1000000	      1386 ns/op
func BenchmarkRBTree_Insert(b *testing.B) {
	rbTree := new(RBTree)
	rbTree.Comparator = container.IntComparator
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = container.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbTree.Insert(data[i])
	}
}
