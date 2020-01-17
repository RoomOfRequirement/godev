package shell

import (
	"godev/basic/algorithm/sort"
)

// Sort implements shell sort
//	https://en.wikipedia.org/wiki/Shellsort
//	changing step size (equivalent to dive array into a table of sub-arrays) + insertion sort
//	step size: n / 2 == n >> 1
func Sort(data sort.Interface) {
	n := data.Len()
	if n < 2 {
		return
	}

	// change step size
	for gap := n >> 1; gap > 0; gap >>= 1 {
		// insertion sort
		for i := gap; i < n; i++ {
			j := i
			for j >= gap && data.Less(j, j-gap) {
				data.Swap(j, j-gap)
				j -= gap
			}
		}
	}
}
