package avltree

import (
	"fmt"
	"goContainer/basic"
	"goContainer/basic/datastructure/tree"
	"testing"
)

func TestNewAVLTree(t *testing.T) {
	var _ basic.Container = (*AVLTree)(nil)
	var _ tree.Tree = (*AVLTree)(nil)

	avlTree := NewAVLTree(basic.IntComparator)

	if !avlTree.Empty() || len(avlTree.Values()) != 0 || len(avlTree.Keys()) != 0 {
		t.Fail()
	}
}

func TestAVLTree_Set(t *testing.T) {
	avlTree := NewAVLTree(basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	aSorted := []int{1, 7, 12, 15, 25, 28, 33, 41}

	for i, k := range a {
		avlTree.Set(k, i)
	}
	if avlTree.Size() != len(a) || len(avlTree.Keys()) != len(a) || len(avlTree.Values()) != len(a) {
		t.Fail()
	}

	for i, v := range aSorted {
		if avlTree.Keys()[i].(int) != v {
			t.Fail()
		}
	}

	for i, v := range avlTree.Values() {
		if a[v.(int)] != aSorted[i] {
			t.Fail()
		}
	}

	for i, k := range a {
		if v, found := avlTree.Get(k); !found || v.(int) != i {
			t.Fail()
		}
	}

	if v, found := avlTree.Get(100); found || v != nil {
		t.Fail()
	}

	avlTree.Clear()

	if !avlTree.Empty() || len(avlTree.Values()) != 0 || len(avlTree.Keys()) != 0 {
		t.Fail()
	}
}

func TestAVLTree_Delete(t *testing.T) {
	avlTree := NewAVLTree(basic.IntComparator)

	a := []int{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		avlTree.Set(a[i], i)
	}

	for i := range a {
		ok := avlTree.Delete(a[i])
		if !ok {
			t.Fail()
		}
	}

	if ok := avlTree.Delete(100); ok {
		t.Fail()
	}

	if !avlTree.Empty() {
		fmt.Println(avlTree.Size())
		t.Fail()
	}
}
