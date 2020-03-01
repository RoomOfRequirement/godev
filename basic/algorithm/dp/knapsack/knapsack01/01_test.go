package knapsack01

import "testing"

func TestZeroOnePack(t *testing.T) {
	items := [][]int{
		{3, 6},
		{5, 4},
		{6, 9},
		{2, 3},
		{7, 8},
	}

	weights := []int{3, 5, 6, 2, 7}
	values := []int{6, 4, 9, 3, 8}

	if res := ZeroOnePack(items, 13); res != 18 {
		t.Fatal("ZeroOnePack: ", res)
	}

	if res := ZeroOnePack1(weights, values, 13); res != 18 {
		t.Fatal("ZeroOnePack1: ", res)
	}

	if res := ZeroOnePackSpaceOpt(items, 13); res != 18 {
		t.Fatal("ZeroOnePackSpaceOpt: ", res)
	}

	if res := ZeroOnePackSpaceOpt1(weights, values, 13); res != 18 {
		t.Fatal("ZeroOnePackSpaceOpt1: ", res)
	}
}
