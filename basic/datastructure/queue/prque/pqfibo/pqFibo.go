package pqfibo

import (
	heap2 "godev/basic/datastructure/heap"
	fibonacci2 "godev/basic/datastructure/heap/fibonacci"
)

// PQFib priority queue based on fibonacci heap
//	https://github.com/Harold2017/godev/tree/master/heap/fibonacci
type PQFib struct {
	items *fibonacci2.Heap
}

// item data struct used in queue
type item struct {
	priority int
	// content to be stored in queue
	content interface{}
}

// Compare to meet heap.item interface
func (it *item) Compare(ait heap2.Item) int {
	if it.priority < ait.(*item).priority {
		return 1
	} else if it.priority == ait.(*item).priority {
		return 0
	} else {
		return -1
	}
}

// newItem to create items with input priority to be stored in the queue
func newItem(x interface{}, priority int) *item {
	return &item{
		priority: priority,
		content:  x,
	}
}

// NewPQFib returns a new priority queue based on fibonacci heap
func NewPQFib() *PQFib {
	return &PQFib{items: fibonacci2.NewHeap()}
}

// Push `x` in item list
func (pq *PQFib) Push(x interface{}, priority int) {
	pq.items.Insert(newItem(x, priority))
}

// Pop content of the item with highest priority
func (pq *PQFib) Pop() interface{} {
	return pq.items.DeleteMin().(*item).content
}

// Size returns queue length
func (pq *PQFib) Size() int {
	return pq.items.Size()
}

// Empty returns true if queue is empty
func (pq *PQFib) Empty() bool {
	return pq.items.Empty()
}

// Clear clears the queue
func (pq *PQFib) Clear() {
	pq.items.Clear()
}

// Values returns values stored int the queue (unordered)
func (pq *PQFib) Values() []interface{} {
	res := make([]interface{}, pq.Size())
	for i, v := range pq.items.Values() {
		res[i] = v.(*item).content
	}
	return res
}
