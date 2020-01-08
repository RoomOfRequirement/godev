package merge

import (
	"goContainer/basic/algorithm/sort"
)

// Sort implements merge sort
//	https://en.wikipedia.org/wiki/Merge_sort
//	notice: this is slice type specified
func Sort(data sort.IntSlice) sort.IntSlice {
	n := data.Len()
	if n < 2 {
		return data
	}
	// divide
	left, right := Sort(data[:n/2]), Sort(data[n/2:])
	// merge
	return merge(left, right)
}

func merge(left, right sort.IntSlice) sort.IntSlice {
	merged := make(sort.IntSlice, left.Len()+right.Len())

	// i for loop left, j for loop right, idx for loop merged
	i, j, idx := 0, 0, 0
	for {
		// `=` for stable sort
		if left[i] <= right[j] {
			merged[idx] = left[i]
			i++
		} else {
			merged[idx] = right[j]
			j++
		}
		idx++
		if i == left.Len() {
			copy(merged[idx:], right[j:])
			break
		}
		if j == right.Len() {
			copy(merged[idx:], left[i:])
			break
		}
	}
	return merged
}
