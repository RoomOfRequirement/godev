package deque

import (
	"fmt"
	"godev/basic"
	"godev/utils"
	"math"
	"strconv"
	"testing"
)

func TestDeque(t *testing.T) {

	var _ basic.Container = (*Deque)(nil)

	dq := NewDeque(-10)

	if dq.Cap() != 8 || !dq.Empty() || dq.Values() != nil {
		t.Fail()
	}

	a := []int{1, 2, 3, 4, 5, 6, 7}
	for _, i := range a {
		dq.PushBack(i)
	}

	if dq.Size() != 7 || dq.PositionsCanPushBack() != 1 || dq.PositionsCanPopFront() != 7 {
		t.Fail()
	}

	for i, v := range dq.Values() {
		if v.(int) != a[i] {
			fmt.Printf("%d", v.(int))
			t.Fail()
		}
	}

	last, err := dq.PopBack()
	if err != nil || last != 7 || dq.Size() != 6 || dq.PositionsCanPushBack() != 2 || dq.PositionsCanPopFront() != 6 {
		t.Fail()
	}

	first, err := dq.PopFront()
	if err != nil || first != 1 || dq.Size() != 5 || dq.PositionsCanPushBack() != 3 || dq.PositionsCanPopFront() != 5 {
		t.Fail()
	}

	dq.PushFront(0)
	f, err := dq.Front()
	if err != nil || f != 0 || dq.Size() != 6 || dq.PositionsCanPushBack() != 2 || dq.PositionsCanPopFront() != 6 {
		t.Fail()
	}

	// [0, 2, 3, 4, 5, 6] + [1, 2, 3, 4, 5, 6, 7]
	for _, i := range a {
		dq.PushBack(i)
	}
	f, _ = dq.Front()
	b, err := dq.Back()
	if dq.Cap() != 16 || dq.Size() != 13 || err != nil || f != 0 || b != 7 || dq.PositionsCanPushBack() != 3 || dq.PositionsCanPopFront() != 13 {
		t.Fail()
	}

	// [7, 6, 5, 4, 3, 2, 1] + [0, 2, 3, 4, 5, 6] + [1, 2, 3, 4, 5, 6, 7]
	for _, i := range a {
		dq.PushFront(i)
	}
	f, _ = dq.Front()
	b, _ = dq.Back()
	if dq.Cap() != 32 || dq.Size() != 20 || f != 7 || b != 7 || dq.PositionsCanPushBack() != 12 || dq.PositionsCanPopFront() != 20 {
		fmt.Println(f, b)
		fmt.Println(dq)
		t.Fail()
	}

	for i, v := range dq.Values() {
		if v.(int) != dq.At(i).(int) {
			t.Fail()
		}
	}

	if dq.At(-1) != 7 || dq.At(-10) != 4 || dq.At(34) != 5 {
		fmt.Println(dq.At(-1).(int), dq.At(-10).(int), dq.At(34).(int))
		t.Fail()
	}

	aReversed := []int{7, 6, 5, 4, 3, 2, 1}
	for i := range aReversed {
		v, err := dq.PopBack()
		if err != nil || v.(int) != aReversed[i] {
			t.Fail()
		}
	}
	f, _ = dq.Front()
	b, _ = dq.Back()
	// size = 13 != 32/4 (8), cap no shrink
	if dq.Cap() != 32 || dq.Size() != 13 || f != 7 || b != 6 || dq.PositionsCanPushBack() != 19 || dq.PositionsCanPopFront() != 13 {
		fmt.Println(dq.Cap(), dq.Size(), f, b)
		t.Fail()
	}

	// [7, 6, 5, 4, 3, 2, 1] + [0, 2, 3, 4, 5, 6]
	for i := range aReversed {
		v, err := dq.PopFront()
		if err != nil || v.(int) != aReversed[i] {
			t.Fail()
		}
	}
	f, _ = dq.Front()
	b, _ = dq.Back()
	// size = 6 < 32/4 (8), cap shrink
	if dq.Cap() != 16 || dq.Size() != 6 || f != 0 || b != 6 || dq.PositionsCanPushBack() != 10 || dq.PositionsCanPopFront() != 6 {
		fmt.Println(dq.Cap(), dq.Size(), f, b)
		t.Fail()
	}

	dq.Clear()

	if !dq.Empty() || dq.Cap() != 16 || dq.Size() != 0 || dq.PositionsCanPushBack() != 16 || dq.PositionsCanPopFront() != 0 {
		t.Fail()
	}

	if f, err = dq.Front(); f != nil || err == nil {
		t.Fail()
	}
	if b, err = dq.Back(); b != nil || err == nil {
		t.Fail()
	}
}

// BenchmarkDeque_PushBack-8   	50000000	        24.2 ns/op
func BenchmarkDeque_PushBack(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	b.ResetTimer()

	dq := NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(data[i])
	}
}

// BenchmarkDeque_PushFront-8   50000000	        26.4 ns/op
func BenchmarkDeque_PushFront(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	b.ResetTimer()

	dq := NewDeque(b.N)
	for i := 0; i < len(data); i++ {
		dq.PushFront(data[i])
	}
}

// BenchmarkDeque_PopFront-8   	100000000	        18.8 ns/op
func BenchmarkDeque_PopFront(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	dq := NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(data[i])
	}
	b.ResetTimer()

	for i := 0; i < len(data); i++ {
		_, _ = dq.PopFront()
	}
}

// BenchmarkDeque_PopBack-8   	100000000	        20.0 ns/op
func BenchmarkDeque_PopBack(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	dq := NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(data[i])
	}
	b.ResetTimer()

	for i := 0; i < len(data); i++ {
		_, _ = dq.PopBack()
	}
}

// BenchmarkDeque_Front-8   	1000000000	        3.00 ns/op
func BenchmarkDeque_Front(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	dq := NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(data[i])
	}
	b.ResetTimer()

	for i := 0; i < len(data); i++ {
		_, _ = dq.Front()
	}
}

// BenchmarkDeque_Back-8   		1000000000	        3.01 ns/op
func BenchmarkDeque_Back(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	dq := NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(data[i])
	}
	b.ResetTimer()

	for i := 0; i < len(data); i++ {
		_, _ = dq.Back()
	}
}

func BenchmarkDeque(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		b.Run("Deque PushBack: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			num := utils.GenerateRandomInt()
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.PopBack()
				dq.PushBack(num)
			}
		})

		b.Run("Deque PushFront: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			num := utils.GenerateRandomInt()
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.PopFront()
				dq.PushFront(num)
			}
		})

		b.Run("Deque PopBack: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				dq.PushBack(rn)
				_, _ = dq.PopBack()
			}
		})

		b.Run("Deque PopFront: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				dq.PushBack(rn)
				_, _ = dq.PopFront()
			}
		})

		b.Run("Deque Front: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.Front()
			}
		})

		b.Run("Deque Back: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = utils.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.Back()
			}
		})
	}
}

/*
BenchmarkDeque/Deque_PushBack:_size-10-8         	50000000	        28.2 ns/op
BenchmarkDeque/Deque_PushBack:_size-100-8        	50000000	        28.1 ns/op
BenchmarkDeque/Deque_PushBack:_size-1000-8       	50000000	        28.2 ns/op

BenchmarkDeque/Deque_PushFront:_size-10-8        	50000000	        27.7 ns/op
BenchmarkDeque/Deque_PushFront:_size-100-8       	50000000	        27.5 ns/op
BenchmarkDeque/Deque_PushFront:_size-1000-8      	50000000	        27.8 ns/op

BenchmarkDeque/Deque_PopBack:_size-10-8          	50000000	        27.8 ns/op
BenchmarkDeque/Deque_PopBack:_size-100-8         	50000000	        28.6 ns/op
BenchmarkDeque/Deque_PopBack:_size-1000-8        	50000000	        28.0 ns/op

BenchmarkDeque/Deque_PopFront:_size-10-8         	50000000	        27.6 ns/op
BenchmarkDeque/Deque_PopFront:_size-100-8        	50000000	        27.4 ns/op
BenchmarkDeque/Deque_PopFront:_size-1000-8       	50000000	        27.6 ns/op

BenchmarkDeque/Deque_Front:_size-10-8            	500000000	        2.94 ns/op
BenchmarkDeque/Deque_Front:_size-100-8           	500000000	        2.95 ns/op
BenchmarkDeque/Deque_Front:_size-1000-8          	500000000	        2.95 ns/op

BenchmarkDeque/Deque_Back:_size-10-8             	500000000	        3.38 ns/op
BenchmarkDeque/Deque_Back:_size-100-8            	500000000	        3.34 ns/op
BenchmarkDeque/Deque_Back:_size-1000-8           	500000000	        3.32 ns/op
*/
