package quick

import (
	"goContainer/basic/algorithm/sort"
)

// Sort implements quick sort in place
//	https://en.wikipedia.org/wiki/Quicksort
//	https://golang.org/src/sort/sort.go
//	divide and conquer
func Sort(data sort.Interface) {
	quickSort(data, 0, data.Len())
}

func quickSort(data sort.Interface, head, tail int) {
	if data.Len() < 2 {
		return
	}
	// minimum sorting unit is 3 numbers, which requires tail - head > 1
	if tail-head <= 1 {
		return
	}

	midLow, midHigh := doPivot(data, head, tail)
	// Avoiding recursion on the larger subproblem guarantees a stack depth of at most lg(tail-head)
	if midLow-head < tail-midHigh {
		quickSort(data, head, midLow)
		quickSort(data, midHigh, tail)
	} else {
		quickSort(data, midHigh, tail)
		quickSort(data, head, midLow)
	}
}

func doPivot(data sort.Interface, head, tail int) (midLow, midHigh int) {
	// minimum sorting unit is 3 numbers
	m := int(uint(head+tail) >> 1) // written like this to avoid integer overflow.
	medianOfThree(data, head, m, tail-1)

	pivot := head
	low, high := head+1, tail-1

	for low < high && data.Less(low, pivot) {
		low++
	}
	mid := low
	for {
		// data[mid] <= pivot
		for mid < high && !data.Less(pivot, mid) {
			mid++
		}

		// data[high-1] < pivot
		for mid < high && data.Less(pivot, high-1) {
			high--
		}

		if mid >= high {
			break
		}

		// data[mid] > pivot; data[high-1] <= pivot
		data.Swap(mid, high-1)
		mid++
		high--
	}

	// swap pivot to mid
	data.Swap(pivot, mid-1)

	return mid - 1, high
}

// data[m0] < data[m1] < data[m2]
func medianOfThree(data sort.Interface, m0, m1, m2 int) {
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}
