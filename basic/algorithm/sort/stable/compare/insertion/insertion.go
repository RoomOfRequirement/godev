package insertion

import (
	"goContainer/basic/algorithm/sort"
)

// Sort implements insertion sort
//	https://en.wikipedia.org/wiki/Insertion_sort
func Sort(data sort.Interface) {
	n := data.Len()
	for i := 1; i < n; i++ {
		for j := i; j > 0; j-- {
			if data.Less(j, j-1) {
				data.Swap(j, j-1)
			} else {
				break
			}
		}
	}
}
