package sort

import "testing"

func TestIntSlice(t *testing.T) {
	is := IntSlice{1, 2, 3}
	if is.Len() != 3 {
		t.Fail()
	}
	if !is.Less(0, 1) || is.Less(1, 0) {
		t.Fail()
	}
	is.Swap(0, 2)
	if is[0] != 3 || is[2] != 1 {
		t.Fail()
	}

	if min, max := is.MinMax(); min != 1 || max != 3 {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	is0 := IntSlice{1, 2}
	is1 := IntSlice{1, 2, 3}
	is2 := IntSlice{1, 2, 3}
	is3 := IntSlice{2, 3, 1}
	if !Equal(is1, is2) || Equal(is2, is3) || Equal(is0, is1) {
		t.Fail()
	}
}
