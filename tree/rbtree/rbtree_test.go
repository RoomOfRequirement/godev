package rbtree

import (
	"fmt"
	"goContainer"
	"goContainer/tree"
	"goContainer/utils"
	"testing"
)

func TestRBTree(t *testing.T) {
	root := NewNode(0, 0)
	rbTree := NewRBTree(root, container.IntComparator)

	var _ tree.Tree = (*RBTree)(nil)

	if rbTree.Empty() {
		t.Fail()
	}

	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for v, k := range a {
		rbTree.Insert(k, v)
	}

	if rbTree.Root.key.(int) != 2 {
		t.Fail()
	}
	if rbTree.MinKey().(int) != -8 {
		t.Fail()
	}

	fmt.Println(rbTree)

	ok := rbTree.Delete(0)
	if !ok {
		t.Fail()
	}

	aSorted := []int{-8, -5, -3, 1, 2, 3, 4, 5, 6, 8}
	for i, v := range rbTree.Keys() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
	expectedValues := map[int]struct{}{
		0: {},
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
		7: {},
		8: {},
		9: {},
	}
	for _, v := range rbTree.Values() {
		if _, found := expectedValues[v.(int)]; !found {
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
	for i, v := range rbTree.Keys() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}

	rbTree.Clear()
	if !rbTree.Empty() || rbTree.Size() != 0 || rbTree.Keys() != nil || rbTree.Values() != nil {
		t.Fail()
	}
}

// BenchmarkRBTree_Insert-8   	 1000000	      1386 ns/op
func BenchmarkRBTree_Insert(b *testing.B) {
	rbTree := new(RBTree)
	rbTree.Comparator = container.IntComparator
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbTree.Insert(data[i], data[i])
	}
}
