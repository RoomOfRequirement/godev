package bheap

import (
	"goContainer"
	"math"
	"strconv"
	"testing"
)

func TestHeap(t *testing.T) {
	var _ container.Container = (*Heap)(nil)
	var _ container.Container = (*MaxHeap)(nil)
	var _ container.Container = (*MinHeap)(nil)
	h := &MaxHeap{
		heap:       nil,
		Comparator: container.IntComparator,
	}
	h.Init()
	a := []int{12, 3, 56, 5, 16, 32, 27, 6, 88}
	for _, v := range a {
		h.Push(v)
	}
	aSorted := []int{88, 56, 32, 16, 5, 12, 27, 3, 6}
	for i, v := range h.heap.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
	if h.Remove(0).(int) != 88 {
		t.Fail()
	}
	if h.Pop().(int) != 56 {
		t.Fail()
	}
	h.Set(3, 102)
	if h.Pop().(int) != 102 {
		t.Fail()
	}

	minH := &MinHeap{
		heap:       nil,
		Comparator: container.IntComparator,
	}
	minH.Init()
	a = []int{12, 3, 56, 5, 16, 32, 27, 6, 88}
	for _, v := range a {
		minH.Push(v)
	}
	aSorted = []int{3, 5, 27, 6, 16, 56, 32, 12, 88}
	for i, v := range minH.heap.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
	if minH.Remove(0).(int) != 3 {
		t.Fail()
	}
	if minH.Pop().(int) != 5 {
		t.Fail()
	}
	minH.Set(3, -2)
	if minH.Pop().(int) != -2 {
		t.Fail()
	}
}

// BenchmarkHeap_Push-8   	10000000	       163 ns/op
func BenchmarkHeap_Push(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = container.GenerateRandomInt()
	}

	minH := new(MinHeap)
	minH.Comparator = container.IntComparator
	minH.Init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		minH.Push(data[i])
	}
}

// BenchmarkHeap_Pop-8   	 1000000	      1706 ns/op
func BenchmarkHeap_Pop(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = container.GenerateRandomInt()
	}

	minH := new(MinHeap)
	minH.Comparator = container.IntComparator
	minH.Init()

	for i := 0; i < b.N; i++ {
		minH.Push(data[i])
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		minH.Pop()
	}
}

func BenchmarkHeap(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		minH := new(MinHeap)
		minH.Comparator = container.IntComparator
		minH.Init()

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()
			minH.Push(rn)
		}

		num := container.GenerateRandomInt()

		b.ResetTimer()

		b.Run("Push: size-"+strconv.Itoa(n), func(b *testing.B) {

			for i := 1; i < b.N; i++ {
				minH.Pop()
				minH.Push(num)
			}
		})

		b.Run("Pop: size-"+strconv.Itoa(n), func(b *testing.B) {

			for i := 1; i < b.N; i++ {
				minH.Push(num)
				minH.Pop()
			}
		})
	}
}

/*
BenchmarkHeap/Push:_size-10-8         	20000000	        99.9 ns/op
BenchmarkHeap/Push:_size-100-8        	20000000	        65.1 ns/op
BenchmarkHeap/Push:_size-1000-8       	20000000	        64.8 ns/op

BenchmarkHeap/Pop:_size-10-8          	20000000	        64.2 ns/op
BenchmarkHeap/Pop:_size-100-8         	20000000	        64.3 ns/op
BenchmarkHeap/Pop:_size-1000-8        	20000000	        64.1 ns/op
*/
