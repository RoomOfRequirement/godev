package bubble

import "goContainer/sort"

// https://en.wikipedia.org/wiki/Bubble_sort

func Sort(data sort.Interface) {
	n := data.Len()
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			if data.Less(j+1, j) {
				data.Swap(j+1, j)
			}
			swapped = true
		}
		if !swapped {
			break
		}
	}
}
