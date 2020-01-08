package heap

import (
	"goContainer/basic/algorithm/sort"
)

// Sort implements heap sort
//	https://en.wikipedia.org/wiki/Heapsort
//	https://golang.org/src/sort/sort.go
func Sort(data sort.Interface) {
	n := data.Len()
	// maximum heap
	for i := (n - 1) >> 1; i >= 0; i-- {
		topDown(data, i, n)
	}

	// pop from root to leaf (max to min) and append into data -> reverse
	for i := n - 1; i >= 0; i-- {
		data.Swap(0, i)
		topDown(data, 0, i)
	}
}

// implement maximum heap (root >= child) on data[low, high] with offset start
func topDown(data sort.Interface, low, high int) {
	root := low
	child := 2*root + 1
	for {
		child = 2*root + 1
		if child >= high {
			break
		}
		if child+1 < high && data.Less(child, child+1) {
			child++
		}
		// root >= child
		if !data.Less(root, child) {
			return
		}
		// swap root with maximum child
		data.Swap(root, child)
		root = child
	}
}
