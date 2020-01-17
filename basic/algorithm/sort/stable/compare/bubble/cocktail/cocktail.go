package cocktail

import (
	"godev/basic/algorithm/sort"
)

// Sort implements cocktail shaker sort
//	https://en.wikipedia.org/wiki/Cocktail_shaker_sort
func Sort(data sort.Interface) {
	n := data.Len()
	left, right := 0, n-1
	for left < right {
		for i := left; i < right; i++ {
			if data.Less(i+1, i) {
				data.Swap(i+1, i)
			}
		}
		right--

		for j := right; j > left; j-- {
			if data.Less(j, j-1) {
				data.Swap(j, j-1)
			}
		}
		left++
	}
}
