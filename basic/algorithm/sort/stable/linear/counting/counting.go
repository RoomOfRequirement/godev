package counting

import (
	"goContainer/basic/algorithm/sort"
)

// Sort implements counting sort, useful for uniformly distributed data over a range
//	https://en.wikipedia.org/wiki/Counting_sort
//	original K range is [0, max+1), here shifts it to [min, max+1), min >= 0 to reduce count array size
//	notice: this is type specified, only for IntSlice here
func Sort(data sort.IntSlice) sort.IntSlice {
	min, max := data.MinMax()
	// max+1 - min instead of max+1 to reduce array size
	count := make(sort.IntSlice, max+1-min)
	sorted := make(sort.IntSlice, data.Len())

	// count every int in data
	// data[i] - min == shift all data to left with min
	for i := 0; i < data.Len(); i++ {
		count[data[i]-min]++
	}

	// update count to count of <= data[x]
	for j := 1; j < count.Len(); j++ {
		count[j] += count[j-1]
	}

	// reverse put
	for k := data.Len() - 1; k >= 0; k-- {
		sorted[count[data[k]-min]-1] = data[k]
		count[data[k]-min]--
	}

	return sorted
}
