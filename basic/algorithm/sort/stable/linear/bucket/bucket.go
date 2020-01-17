package bucket

import (
	"godev/basic/algorithm/sort"
	"godev/basic/algorithm/sort/stable/compare/insertion"
)

// Sort implements bucket (bin) sort, useful for uniformly distributed data over a range
//	https://en.wikipedia.org/wiki/Bucket_sort
//	notice: this is type specified, only for IntSlice here
//	how to efficiently decide bucket (step) size?
func Sort(data sort.IntSlice) {
	// according to <Introduction To Algorithm>,
	// square of bucket size should be proportional to elements number
	// to achieve O(n) time complexity

	// if bucketSize = 1, it will become counting sort (because counting is only for int and int "width" is 1)
	// bucket sort breaks continuous data into discrete data
	// for other bucketSize, its complexity is proportional to sort algorithm complexity used in bucket
	bucketSize := 10 //int(math.Sqrt(float64(data.Len())))
	min, max := data.MinMax()
	bucketNum := (max-min)/bucketSize + 1
	buckets := make([]sort.IntSlice, bucketNum)
	for i := range buckets {
		buckets[i] = sort.IntSlice{}
	}

	for i := 0; i < data.Len(); i++ {
		idx := (data[i] - min) / bucketSize
		buckets[idx] = append(buckets[idx], data[i])
	}

	idx := 0
	for i := range buckets {
		insertion.Sort(buckets[i])
		for _, k := range buckets[i] {
			data[idx] = k
			idx++
		}
	}
}
