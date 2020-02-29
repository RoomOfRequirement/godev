package quickselect

import (
	"math/rand"
	"sort"
)

// QuickSelect to find kth smallest element in slice
//	notice: k starts from 0
// https://rosettacode.org/wiki/Quickselect_algorithm#Go
func QuickSelect(arr sort.Interface, k int) (idx int) {
	first := 0
	last := arr.Len() - 1
	for {
		pivotIdx := partition(arr, first, last, rand.Intn(last-first+1)+first)
		if k == pivotIdx {
			return pivotIdx
		} else if k < pivotIdx {
			last = pivotIdx - 1
		} else {
			first = pivotIdx + 1
		}
	}
}

func partition(arr sort.Interface, first, last, pivotIdx int) int {
	arr.Swap(first, pivotIdx) // move it to beginning
	left := first + 1
	right := last
	for left <= right {
		for left <= last && arr.Less(left, first) {
			left++
		}
		for right >= first && arr.Less(first, right) {
			right--
		}
		if left <= right {
			arr.Swap(left, right)
			left++
			right--
		}
	}
	arr.Swap(first, right)
	return right
}

// MOM based on median of medians
//	notice: k starts from 0
func MOM(arrInt []int, k int) int {
	// k + 1 due to k starts from 0
	return medianOfMedians(arrInt, k+1, 5)
}

func medianOfMedians(arrInt []int, k, r int) int {

	num := len(arrInt)
	if num < 10 {
		sort.Ints(arrInt)
		return arrInt[k-1]
	}
	med := (num + r - 1) / r

	medians := make([]int, med)

	for i := 0; i < med; i++ {
		v := (i * r) + r
		var arr []int
		if v >= num {
			arr = make([]int, len(arrInt[(i*r):]))
			copy(arr, arrInt[(i*r):])
		} else {
			arr = make([]int, r)
			copy(arr, arrInt[(i*r):v])
		}
		sort.Ints(arr)
		medians[i] = arr[len(arr)/2]
	}
	pivot := medianOfMedians(medians, (len(medians)+1)/2, r)

	var leftSide, rightSide []int

	for i := range arrInt {
		if arrInt[i] < pivot {
			leftSide = append(leftSide, arrInt[i])
		} else if arrInt[i] > pivot {
			rightSide = append(rightSide, arrInt[i])
		}
	}

	switch {
	case k == (len(leftSide) + 1):
		return pivot
	case k <= len(leftSide):
		return medianOfMedians(leftSide, k, r)
	default:
		return medianOfMedians(rightSide, k-len(leftSide)-1, r)
	}
}
