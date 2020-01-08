package knapsack01

import "goContainer/basic/datastructure/bits"

// 0-1 knapsack problem
//	https://en.wikipedia.org/wiki/Knapsack_problem
//	state transfer function: f[i][V] = max{f[i-1][V], f[i-1][V-w[i]] + v[i]}
//	f: total value function
//	V: total volume
//	w[i]: i-th item weight (volume cost)
//	v[i]: i-th item value
//	to i-th item, it can be put in bag (1) or not (0)
//	if not put in (0): f[i][V] = f[i-1][V]; if put in (1): f[i][V] = f[i-1][V-w[i]] + v[i])

// ZeroOnePack returns max value of items with a limited capacity
//	items row: i-th item (items[i])
//	items column: weight (items[i][0]), value (items[i][1])
//	f[i][j] = max(f[i-1][j], f[i-1][j-items[i][0]] + items[i][1])
func ZeroOnePack(items [][]int, capacity int) int {
	itemNum := len(items)
	f := make([][]int, itemNum)
	for i := 0; i < itemNum; i++ {
		f[i] = make([]int, capacity+1) // f[i][0] -> f[i][capacity]
	}

	// f[0]
	for j := items[0][0]; j <= capacity; j++ {
		f[0][j] = items[0][1]
	}

	for i := 1; i < itemNum; i++ {
		for j := 0; j <= capacity; j++ {
			if j >= items[i][0] {
				f[i][j] = max(f[i-1][j], f[i-1][j-items[i][0]]+items[i][1])
			} else {
				f[i][j] = f[i-1][j]
			}
		}
	}

	return f[itemNum-1][capacity]
}

// ZeroOnePack1 same with former but input are divided into two 1D array
func ZeroOnePack1(weights []int, values []int, capacity int) int {
	itemNum := len(weights)

	f := make([][]int, itemNum)
	for i := 0; i < itemNum; i++ {
		f[i] = make([]int, capacity+1) // f[i][0] -> f[i][capacity]
	}

	// f[0]
	for j := weights[0]; j <= capacity; j++ {
		f[0][j] = values[0]
	}

	for i := 1; i < itemNum; i++ {
		for j := 0; j <= capacity; j++ {
			if j >= weights[i] {
				f[i][j] = max(f[i-1][j], f[i-1][j-weights[i]]+values[i])
			} else {
				f[i][j] = f[i-1][j]
			}
		}
	}

	return f[itemNum-1][capacity]
}

// ZeroOnePackSpaceOpt optimizes space consumption by using 1D array {O(V): total volume V} instead of 2D array {O(nV): n items, total volume V} to represent f
//	time consumption is still O(nV)
func ZeroOnePackSpaceOpt(items [][]int, capacity int) int {
	itemNum := len(items)
	f := make([]int, capacity+1)
	for j := items[0][0]; j <= capacity; j++ {
		f[j] = items[0][1]
	}

	for i := 1; i < itemNum; i++ {
		for j := capacity; j >= 0; j-- {
			if j >= items[i][0] {
				f[j] = max(f[j], f[j-items[i][0]]+items[i][1])
			}
		}
	}
	return f[capacity]
}

// ZeroOnePackSpaceOpt1 same with former but input are divided into two 1D array
func ZeroOnePackSpaceOpt1(weights []int, values []int, capacity int) int {
	itemNum := len(weights)

	f := make([]int, capacity+1)
	for j := weights[0]; j <= capacity; j++ {
		f[j] = values[0]
	}

	for i := 1; i < itemNum; i++ {
		for j := capacity; j >= 0; j-- {
			if j >= weights[i] {
				f[j] = max(f[j], f[j-weights[i]]+values[i])
			}
		}
	}
	return f[capacity]
}

func max(a, b int) int {
	return bits.Max(a, b)
}
