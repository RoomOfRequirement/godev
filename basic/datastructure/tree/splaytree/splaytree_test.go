package splaytree

import (
	"godev/basic"
	"godev/basic/datastructure/tree"
	"godev/utils"
	"testing"
)

func TestNewSplayTree(t *testing.T) {
	var _ tree.Tree = (*SplayTree)(nil)

	st := NewSplayTree(basic.IntComparator)
	if !st.Empty() {
		t.Fail()
	}
}

func TestSplayTree_Insert(t *testing.T) {
	a := make([]int, 10)
	for i := 0; i < 10; i++ {
		a[i] = utils.GenerateRandomInt()
	}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	if st.Size() != len(a) {
		t.Fail()
	}
	err := st.Insert(1, 100)
	if err == nil {
		t.Fail()
	}
	expectedErr := KeyAlreadyExistError{1}
	if err.Error() != expectedErr.Error() {
		t.Fail()
	}
}

func TestSplayTree_Get(t *testing.T) {
	a := make([]int, 10)
	for i := 0; i < 10; i++ {
		a[i] = utils.GenerateRandomInt()
	}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, v := range a {
		if st.Get(i).(int) != v {
			t.Fail()
		}
	}
}

func TestSplayTree_Search(t *testing.T) {
	a := []int{1, 8, 5, 3, 7, 2, 6}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := range a {
		if !st.Search(i) {
			t.Fail()
		}
	}

	if st.Search(10) {
		t.Fail()
	}
}

func TestSplayTree_Delete(t *testing.T) {
	a := make([]int, 10)
	for i := 0; i < 10; i++ {
		a[i] = utils.GenerateRandomInt()
	}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := range a {
		if err := st.Delete(i); err != nil {
			t.Fatal(err)
		}
	}

	if !st.Empty() {
		t.Fail()
	}

	err := st.Insert(1, 100)
	if err != nil {
		t.Fatal(err)
	}

	err = st.Delete(10)
	if err == nil {
		t.Fail()
	}
	expectedErr := KeyNotExistError{10}
	if err.Error() != expectedErr.Error() {
		t.Fail()
	}
}

func TestSplayTree_Update(t *testing.T) {
	a := []int{1, 8, 5, 3, 7, 2, 6}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	err := st.Update(2, 10)
	if err != nil {
		t.Fatal(err)
	}
	if st.Get(2).(int) != 10 {
		t.Fail()
	}

	err = st.Update(10, 100)
	if err == nil {
		t.Fail()
	}
}

func TestSplayTree_Clear(t *testing.T) {
	a := []int{1, 8, 5, 3, 7, 2, 6}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	st.Clear()
	if !st.Empty() {
		t.Fail()
	}
}

func TestSplayTree_Keys(t *testing.T) {
	a := []int{1, 8, 5, 3, 7, 2, 6}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	expected := []int{6, 5, 4, 3, 2, 1, 0}

	for i, v := range st.Keys() {
		if v.(int) != expected[i] {
			t.Fail()
		}
	}
}

func TestSplayTree_Values(t *testing.T) {
	a := []int{1, 8, 5, 3, 7, 2, 6}
	st := NewSplayTree(basic.IntComparator)
	for i, v := range a {
		err := st.Insert(i, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	expected := []int{6, 5, 4, 3, 2, 1, 0}

	for i, v := range st.Values() {
		if v.(int) != a[expected[i]] {
			t.Fail()
		}
	}
}

func TestKeyAlreadyExistError_Error(t *testing.T) {
	err := KeyAlreadyExistError{key: 1}
	if err.key.(int) != 1 || err.Error() != "insert error: key 1 already exists, you can choose Update" {
		t.Fail()
	}
}

func TestKeyNotExistError_Error(t *testing.T) {
	err := KeyNotExistError{key: 1}
	if err.key.(int) != 1 || err.Error() != "delete error: key 1 does NOT exist" {
		t.Fail()
	}
}

// BenchmarkSplayTree_Insert-4   	 1000000	      2609 ns/op
func BenchmarkSplayTree_Insert(b *testing.B) {
	st := NewSplayTree(basic.IntComparator)
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = st.Insert(data[i], data[i])
	}
}

// BenchmarkSplayTree_Delete-4   	 1000000	      2184 ns/op
func BenchmarkSplayTree_Delete(b *testing.B) {
	st := NewSplayTree(basic.IntComparator)
	data := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = utils.GenerateRandomInt()
	}

	for i := 0; i < b.N; i++ {
		_ = st.Insert(data[i], data[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = st.Delete(data[i])
	}
}
