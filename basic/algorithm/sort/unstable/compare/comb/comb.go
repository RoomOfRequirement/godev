package comb

import (
	"godev/basic/algorithm/sort"
)

// Sort implements comb sort
//	https://en.wikipedia.org/wiki/Comb_sort
//	similar with Shell sort (improved insertion sort) but for bubble sort
func Sort(data sort.Interface) {
	n := data.Len()
	shrinkFactor := 0.8
	gap := n
	swapped := true

	for gap > 1 || swapped {
		if gap > 1 {
			gap = int(float64(gap) * shrinkFactor)
		}

		// bubble sort
		swapped = false
		for i := 0; i < n-gap; i++ {
			if data.Less(i+gap, i) {
				data.Swap(i+gap, i)
				swapped = true
			}
		}
	}
}
