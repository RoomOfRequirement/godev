package concurrent

import (
	"fmt"
	"goContainer"
	"goContainer/queue/deque"
	"sync"
	"testing"
)

func TestNewDequeRW(t *testing.T) {
	// meet container interface
	var _ container.Container = (*DequeRW)(nil)

	dq := NewDequeRW(-10)
	if dq.Cap() < 0 || dq.Size() != 0 || !dq.Empty() || dq.Values() != nil || dq.PositionsCanPopFront() != 0 || dq.PositionsCanPushBack() < 0 {
		t.Fail()
	}
}

func TestDequeRW_Cap(t *testing.T) {
	dq := NewDequeRW(100) // 128
	if dq.Cap() != 128 || dq.PositionsCanPushBack() != 128 {
		t.Fail()
	}
}

func TestDequeRW_At_Front_Back(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	for i := range a {
		dq.PushBack(i)
	}
	if dq.At(0) != 0 || dq.At(1) != 1 || dq.At(2) != 2 {
		t.Fail()
	}

	f, err := dq.Front()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	b, err := dq.Back()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if f != 0 || b != 2 {
		t.Fail()
	}

	dq.Clear()
	if !dq.Empty() {
		t.Fail()
	}
}

func TestDequeRW_IsFull(t *testing.T) {
	dq := NewDequeRW(20) // 32
	for i := 0; i < 20; i++ {
		dq.PushBack(i)
	}
	if dq.IsFull() || dq.Size() != 20 {
		t.Fail()
	}
	for i := 0; i < 12; i++ {
		dq.PushBack(i)
	}
	if !dq.IsFull() || dq.Size() != 32 {
		t.Fail()
	}
}

func TestDequeRW_CC_PB(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	wg := sync.WaitGroup{}
	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushBack(i)
		}(dq, i)
	}
	wg.Wait()
	if dq.Size() != 3 || dq.PositionsCanPushBack() != 29 || dq.PositionsCanPopFront() != 3 {
		fmt.Println(dq)
		t.Fail()
	}
}

func TestDequeRW_CC_PF(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	wg := sync.WaitGroup{}
	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushFront(i)
		}(dq, i)
	}
	wg.Wait()

	if dq.Size() != 3 || dq.PositionsCanPushBack() != 29 || dq.PositionsCanPopFront() != 3 {
		fmt.Println(dq)
		t.Fail()
	}
}

func TestDequeRW_CC_PB_PF(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	wg := sync.WaitGroup{}
	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushBack(i)
		}(dq, i)

		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushFront(i)
		}(dq, i)
	}
	wg.Wait()

	if dq.Size() != 6 || dq.PositionsCanPushBack() != 26 || dq.PositionsCanPopFront() != 6 {
		fmt.Println(dq)
		t.Fail()
	}
}

func TestDequeRW_CC_PPPB(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	wg := sync.WaitGroup{}
	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushBack(i)
		}(dq, i)
	}
	wg.Wait()

	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			_, err := dq.PopFront()
			if err != nil {
				fmt.Println(err)
				t.Fail()
			}
		}(dq, i)
	}
	wg.Wait()

	if dq.Size() != 0 || dq.PositionsCanPushBack() != 32 || dq.PositionsCanPopFront() != 0 {
		fmt.Println(dq)
		t.Fail()
	}
}

func TestDequeRW_CC_PPPF(t *testing.T) {
	dq := NewDequeRW(20) // 32
	a := []int{0, 1, 2}
	wg := sync.WaitGroup{}
	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			dq.PushBack(i)
		}(dq, i)
	}
	wg.Wait()

	for i := range a {
		wg.Add(1)
		go func(dq *DequeRW, i int) {
			defer wg.Done()
			_, err := dq.PopBack()
			if err != nil {
				fmt.Println(err)
				t.Fail()
			}
		}(dq, i)
	}
	wg.Wait()

	if dq.Size() != 0 || dq.PositionsCanPushBack() != 32 || dq.PositionsCanPopFront() != 0 {
		fmt.Println(dq)
		t.Fail()
	}
}

// BenchmarkDequeRW_PushBack-8      5000000               232 ns/op               0 B/op          0 allocs/op
func BenchmarkDequeRW_PushBack(b *testing.B) {
	dq := NewDequeRW(b.N)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dq.PushBack(0)
		}
	})
}

// BenchmarkDeque_PushBack-8       300000000                5.38 ns/op            0 B/op          0 allocs/op
func BenchmarkDeque_PushBack(b *testing.B) {
	dq := deque.NewDeque(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dq.PushBack(0)
	}
}

// BenchmarkDequeRW_PushBack_100Elem-8        50000             28372 ns/op            4856 B/op         99 allocs/op
func BenchmarkDequeRW_PushBack_100Elem(b *testing.B) {
	dq := NewDequeRW(b.N)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 100; i++ {
				dq.PushBack(i)
			}
		}
	})
}

// BenchmarkDeque_PushBack_100Elem-8         500000              4888 ns/op            4856 B/op         99 allocs/op
func BenchmarkDeque_PushBack_100Elem(b *testing.B) {
	dq := deque.NewDeque(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			dq.PushBack(i)
		}
	}
}

// BenchmarkDequeRW_PopBack-8       5000000               259 ns/op              15 B/op          0 allocs/op
func BenchmarkDequeRW_PopBack(b *testing.B) {
	dq := NewDequeRW(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(i)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = dq.PopBack()
		}
	})
}

// BenchmarkDeque_PopBack-8        100000000               21.2 ns/op            15 B/op          0 allocs/op
func BenchmarkDeque_PopBack(b *testing.B) {
	dq := deque.NewDeque(b.N)
	for i := 0; i < b.N; i++ {
		dq.PushBack(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = dq.PopBack()
	}
}

// BenchmarkDequeRW_PopBack_100Elem-8         30000             44834 ns/op             437 B/op         53 allocs/op
func BenchmarkDequeRW_PopBack_100Elem(b *testing.B) {
	dq := NewDequeRW(b.N)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			for i := 0; i < 100; i++ {
				dq.PushBack(i)
			}

			b.StartTimer()
			for i := 0; i < 100; i++ {
				_, _ = dq.PopBack()
			}
		}
	})
}

// BenchmarkDeque_PopBack_100Elem-8         2000000               759 ns/op               0 B/op          0 allocs/op
func BenchmarkDeque_PopBack_100Elem(b *testing.B) {
	dq := deque.NewDeque(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for j := 0; j < 100; j++ {
			dq.PushBack(j)
		}

		b.StartTimer()
		for j := 0; j < 100; j++ {
			_, _ = dq.PopBack()
		}
	}
}
