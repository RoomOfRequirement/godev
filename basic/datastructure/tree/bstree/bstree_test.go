package bstree

import (
	"goContainer/basic"
	"goContainer/basic/datastructure/tree"
	"goContainer/utils"
	"testing"
)

func TestBSTree_InterfaceAssertion(t *testing.T) {
	var _ tree.Tree = (*BSTree)(nil)
}

func TestBSTree_Insert(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	aSorted := []int{-8, -5, -3, 1, 2, 3, 4, 5, 6, 8}
	for i, v := range bst.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_Empty(t *testing.T) {
	bst := NewBSTree()
	if !bst.Empty() {
		t.Fail()
	}
	bst.Insert(1)
	if bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Size(t *testing.T) {
	bst := NewBSTree()
	bst.Insert(1)
	if bst.Size() != 1 {
		t.Fail()
	}
}

func TestBSTree_Clear(t *testing.T) {
	bst := NewBSTree()
	bst.Insert(1)
	bst.Clear()
	if !bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Lookup(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	for i := range a {
		if !bst.Lookup(a[i]) {
			t.Fail()
		}
	}
	if bst.Lookup(10) {
		t.Fail()
	}
}

func TestBSTree_MaxDepth(t *testing.T) {
	bst := NewBSTree()
	if bst.MaxDepth() != 0 {
		t.Fail()
	}
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if bst.MaxDepth() != 6 {
		t.Fail()
	}
}

func TestBSTree_MinValue(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if bst.MinValue() != -8 {
		t.Fail()
	}
}

func TestBSTree_HasPathSum(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	bst.Diffidence = func(sum, data interface{}) interface{} {
		return sum.(int) - data.(int)
	}
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if !bst.HasPathSum(-12) || !bst.HasPathSum(-7) || !bst.HasPathSum(15) || !bst.HasPathSum(25) {
		t.Fail()
	}
}

func TestBSTree_PrintPaths(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.PrintPaths()
}

func TestBSTree_Mirror(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.Mirror()
	aReverseSorted := []int{8, 6, 5, 4, 3, 2, 1, -3, -5, -8}
	for i, v := range bst.Values() {
		if v.(int) != aReverseSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_DoubleTree(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.DoubleTree()
	aSorted := []int{1, 1, 2, 2, 3, 3}
	for i, v := range bst.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_SameTree(t *testing.T) {
	bstA, bstB := NewBSTree(), NewBSTree()
	bstA.Comparator, bstB.Comparator = basic.IntComparator, basic.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bstA.Insert(v)
		bstB.Insert(v)
	}
	if !bstA.SameTree(bstB) {
		t.Fail()
	}
	bstB.Insert(5)
	if bstA.SameTree(bstB) {
		t.Fail()
	}
}

func TestCountTrees(t *testing.T) {
	if CountTrees(4) != 14 {
		t.Fail()
	}
}

func TestBSTree_IsBST(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = basic.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bst.Insert(v)
	}
	if !bst.IsBST(0, 5) {
		t.Fail()
	}
	bst.Clear()
	bst.Root = NewNode(1)
	bst.Root.left = NewNode(3)
	bst.Root.right = NewNode(2)
	if bst.IsBST(0, 5) {
		t.Fail()
	}
}

// BenchmarkBSTree_Insert-8   	 1000000	      1353 ns/op
func BenchmarkBSTree_Insert(b *testing.B) {
	bst := new(BSTree)
	bst.Comparator = basic.IntComparator
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.Insert(data[i])
	}
}
