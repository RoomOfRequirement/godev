package radix

import (
	"godev/basic/algorithm/sort"
	"math"
)

// Sort implements radix sort, useful for uniformly distributed data over a range
//	https://en.wikipedia.org/wiki/Radix_sort
//	notice: this is type specified, only for IntSlice here
//	time complexity O(n * k), k ~= logB N, B is radix, for int, B = 10
func Sort(data sort.IntSlice) {
	_, max := data.MinMax()
	radix := 10
	// K digits number
	K := int(math.Ceil(math.Log10(float64(max + 1))))
	for i := 1; i < K+1; i++ {
		// radix buckets
		buckets := make([]sort.IntSlice, radix)
		for i := range buckets {
			buckets[i] = sort.IntSlice{}
		}
		for _, num := range data {
			// K-th digit, from low to high (right to left)
			bucketIdx := num % int(math.Pow(float64(radix), float64(i))) / int(math.Pow(float64(radix), float64(i-1)))
			buckets[bucketIdx] = append(buckets[bucketIdx], num)
		}

		// merge buckets to update data
		idx := 0
		for i := range buckets {
			for _, k := range buckets[i] {
				data[idx] = k
				idx++
			}
		}
	}
}
