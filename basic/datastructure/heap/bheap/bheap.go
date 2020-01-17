package bheap

import (
	"container/heap"
	"godev/basic"
)

// Heap use heap interface from "container/heap"
type Heap struct {
	values     []interface{}
	comparator basic.Comparator
}

func (h Heap) Len() int {
	return len(h.values)
}

func (h Heap) Less(i, j int) bool {
	return h.comparator(h.values[i], h.values[j]) < 0
}

func (h Heap) Swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

// Push pushes value into the heap
func (h *Heap) Push(v interface{}) {
	h.values = append(h.values, v)
}

// Pop pops heap roof
func (h *Heap) Pop() interface{} {
	size := len(h.values)
	lastNode := h.values[size-1]
	// avoid memory leak
	h.values[size-1] = nil
	h.values = h.values[:size-1]
	return lastNode
}

// Size returns the number of values inside the heap
func (h *Heap) Size() int {
	return len(h.values)
}

// Empty returns true if no value inside the heap
func (h *Heap) Empty() bool {
	return len(h.values) == 0
}

// Clear clears values inside the heap
func (h *Heap) Clear() {
	h.values = ([]interface{})(nil)
}

// Values returns values inside the heap
func (h *Heap) Values() []interface{} {
	return h.values
}

func (h *Heap) set(i int, v interface{}) {
	if i > len(h.values)-1 || i < 0 {
		panic("invalid index")
	}
	h.values[i] = v
}

// MinHeap heap stored minimum value in its root
type MinHeap struct {
	heap       *Heap
	Comparator basic.Comparator
}

// Init initializes MinHeap
func (h *MinHeap) Init() {
	h.heap = &Heap{
		values:     nil,
		comparator: h.Comparator,
	}
	heap.Init(h.heap)
}

// Push pushes values into the heap
func (h *MinHeap) Push(v interface{}) {
	heap.Push(h.heap, v)
}

// Pop pops heap roof
func (h *MinHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

// Remove removes value with index `i`
func (h *MinHeap) Remove(i int) interface{} {
	return heap.Remove(h.heap, i)
}

// Set sets index `i` with value `v`
func (h *MinHeap) Set(i int, v interface{}) {
	h.heap.set(i, v)
	heap.Fix(h.heap, i)
}

// Size returns the number of values inside the heap
func (h *MinHeap) Size() int {
	return len(h.heap.values)
}

// Empty returns true if no value inside the heap
func (h *MinHeap) Empty() bool {
	return len(h.heap.values) == 0
}

// Clear clears values inside the heap
func (h *MinHeap) Clear() {
	h.heap.values = ([]interface{})(nil)
}

// Values returns values inside the heap
func (h *MinHeap) Values() []interface{} {
	return h.heap.values
}

// MaxHeap heap stored maximum value in its root
type MaxHeap struct {
	heap       *Heap
	Comparator basic.Comparator
}

// Init initializes MaxHeap
func (h *MaxHeap) Init() {
	h.heap = &Heap{
		values: nil,
		comparator: func(a, b interface{}) int {
			return -h.Comparator(a, b)
		},
	}
	heap.Init(h.heap)
}

// Push pushes values into the heap
func (h *MaxHeap) Push(v interface{}) {
	heap.Push(h.heap, v)
}

// Pop pops heap roof
func (h *MaxHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

// Remove removes value with index `i`
func (h *MaxHeap) Remove(i int) interface{} {
	return heap.Remove(h.heap, i)
}

// Set sets index `i` with value `v`
func (h *MaxHeap) Set(i int, v interface{}) {
	h.heap.set(i, v)
	heap.Fix(h.heap, i)
}

// Size returns the number of values inside the heap
func (h *MaxHeap) Size() int {
	return len(h.heap.values)
}

// Empty returns true if no value inside the heap
func (h *MaxHeap) Empty() bool {
	return len(h.heap.values) == 0
}

// Clear clears values inside the heap
func (h *MaxHeap) Clear() {
	h.heap.values = ([]interface{})(nil)
}

// Values returns values inside the heap
func (h *MaxHeap) Values() []interface{} {
	return h.heap.values
}
