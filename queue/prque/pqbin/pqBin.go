package pqbin

import (
	"goContainer/heap/bheap"
)

// PQBin priority queue based on binary heap, implements golang heap.Interface (Sort.Interface)
//	https://github.com/Harold2017/golina/tree/master/container/heap/bheap
type PQBin struct {
	items *bheap.MaxHeap
}

// item data struct used in queue
type item struct {
	priority int
	// content to be stored in queue
	content interface{}
}

// newItem to create items with input priority to be stored in the queue
func newItem(x interface{}, priority int) *item {
	return &item{
		priority: priority,
		content:  x,
	}
}

func itemComparator(a, b interface{}) int {
	itemA, itemB := a.(*item), b.(*item)
	A, B := itemA.priority, itemB.priority
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

// Push `x` in item list
//	 O(log n)
func (pq *PQBin) Push(x interface{}, priority int) {
	pq.items.Push(newItem(x, priority))
}

// Pop content of the item with highest priority
//	 O(log n)
func (pq *PQBin) Pop() interface{} {
	return pq.items.Pop().(*item).content
}

// NewPQBin returns a new priority queue based on golang Heap interface
func NewPQBin() *PQBin {
	maxH := bheap.MaxHeap{
		Comparator: itemComparator,
	}
	maxH.Init()
	return &PQBin{items: &maxH}
}

// Size returns queue length
func (pq *PQBin) Size() int {
	return pq.items.Size()
}

// Empty returns true if queue is empty
func (pq *PQBin) Empty() bool {
	return pq.items.Empty()
}

// Clear clears the queue
func (pq *PQBin) Clear() {
	pq.items.Clear()
}

// Values returns values stored int the queue (unordered)
func (pq *PQBin) Values() []interface{} {
	res := make([]interface{}, pq.Size())
	for i, v := range pq.items.Values() {
		res[i] = v.(*item).content
	}
	return res
}
