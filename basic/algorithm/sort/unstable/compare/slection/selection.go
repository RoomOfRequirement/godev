package slection

import (
	"godev/basic/algorithm/sort"
)

// Sort implements selection sort
//	https://en.wikipedia.org/wiki/Selection_sort
func Sort(data sort.Interface) {
	n := data.Len()

	for i := 0; i < n-1; i++ {
		min := i
		// unsorted elements
		for j := i + 1; j < n; j++ {
			if data.Less(j, min) {
				min = j
			}
		}

		// exchange
		// if min != i
		data.Swap(i, min)
	}
}
