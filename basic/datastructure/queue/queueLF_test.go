package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	c := uint32(8)
	q := NewLockFree(c)
	assert.Equal(t, c, q.Cap())
	assert.True(t, q.Empty())
	ok, size := q.Enqueue("test")
	assert.True(t, ok)
	assert.Equal(t, uint32(1), size)
	assert.False(t, q.Empty())

	ok, val, size := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "test", val.(string))
	assert.Equal(t, uint32(0), size)
	assert.Equal(t, uint32(0), q.Size())
	assert.True(t, q.Empty())

	q = NewLockFree(15)
	assert.Equal(t, uint32(16), q.Cap())
	assert.Equal(t, fmt.Sprintf("Queue: \n\tCap: %d\n\tSize: %d\n\tBuffer: %+v\n",
		q.cap, q.Size(), q.buf), q.String())
	for i := 0; i < 16; i++ {
		q.Enqueue(i)
	}
	assert.Equal(t, uint32(15), q.Size())
	cnt := 0
	for {
		if ok, val, size := q.Dequeue(); ok {
			assert.Equal(t, cnt, val)
			assert.Equal(t, uint32(14-cnt), size)
			cnt++
		} else {
			break
		}
	}
}

// tests from https://github.com/yireyun/go-queue
func TestQueueEnqueueDequeue(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cnt := 10000
	sum := 0
	start := time.Now()
	var EnqueueD, DequeueD time.Duration
	for i := 0; i <= runtime.NumCPU()*4; i++ {
		sum += i * cnt
		Enqueue, Dequeue := testQueueEnqueueDequeue(t, i, cnt)
		EnqueueD += Enqueue
		DequeueD += Dequeue
	}
	end := time.Now()
	use := end.Sub(start)
	op := use / time.Duration(sum)
	t.Logf("Grp: %d, Times: %d, use: %v, %v/op", runtime.NumCPU()*4, sum, use, op)
	t.Logf("Enqueue: %d, use: %v, %v/op", sum, EnqueueD, EnqueueD/time.Duration(sum))
	t.Logf("Dequeue: %d, use: %v, %v/op", sum, DequeueD, DequeueD/time.Duration(sum))
}

func TestQueueGeneral(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		cnt := 10000 * 10
		if i > 9 {
			cnt = 10000 * 1
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueueGeneral(t, i, cnt)
		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func TestQueueEnqueueGoDequeue(t *testing.T) {
	var Sum, miss int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		cnt := 10000 * 10
		if i > 9 {
			cnt = 10000 * 1
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueueEnqueueGoDequeue(t, i, cnt)

		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func TestQueueEnqueueDoDequeue(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		cnt := 10000 * 10
		if i > 9 {
			cnt = 10000 * 1
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueueEnqueueDoDequeue(t, i, cnt)
		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func testQueueEnqueueDequeue(t *testing.T, grp, cnt int) (
	Enqueue time.Duration, Dequeue time.Duration) {
	var wg sync.WaitGroup
	var id int32
	wg.Add(grp)
	q := NewLockFree(1024 * 1024)
	start := time.Now()
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				val := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&id, 1))
				ok, _ := q.Enqueue(&val)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Enqueue(&val)
				}
			}
		}(i)
	}
	wg.Wait()
	end := time.Now()
	Enqueue = end.Sub(start)

	wg.Add(grp)
	start = time.Now()
	for i := 0; i < grp; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < cnt; {
				ok, _, _ := q.Dequeue()
				if !ok {
					runtime.Gosched()
				} else {
					j++
				}
			}
		}()
	}
	wg.Wait()
	end = time.Now()
	Dequeue = end.Sub(start)
	if q := q.Size(); q != 0 {
		t.Errorf("Grp:%v, Size Error: [%v] <>[%v]", grp, q, 0)
	}
	return Enqueue, Dequeue
}

func testQueueGeneral(t *testing.T, grp, cnt int) int {

	var wg sync.WaitGroup
	var idPut, idGet int32
	var miss int32

	wg.Add(grp)
	q := NewLockFree(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				val := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&idPut, 1))
				ok, _ := q.Enqueue(&val)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Enqueue(&val)
				}
			}
		}(i)
	}

	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _, _ = q.Dequeue()
				for !ok {
					atomic.AddInt32(&miss, 1)
					time.Sleep(time.Microsecond * 50)
					ok, _, _ = q.Dequeue()
				}
				atomic.AddInt32(&idGet, 1)
			}
		}(i)
	}
	wg.Wait()
	if q := q.Size(); q != 0 {
		t.Errorf("Grp:%v, Size Error: [%v] <>[%v]", grp, q, 0)
	}
	return int(miss)
}

var value = 1

func testQueueEnqueueGoDequeue(t *testing.T, grp, cnt int) int {
	var wg sync.WaitGroup
	wg.Add(grp)
	q := NewLockFree(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Enqueue(&value)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Enqueue(&value)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _, _ = q.Dequeue()
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _, _ = q.Dequeue()
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0
}

func testQueueEnqueueDoDequeue(t *testing.T, grp, cnt int) int {
	var wg sync.WaitGroup
	wg.Add(grp)
	q := NewLockFree(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Enqueue(&value)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Enqueue(&value)
				}
				ok, _, _ = q.Dequeue()
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _, _ = q.Dequeue()
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0
}

func testQueueEnqueueDequeueOrder(t *testing.T, grp, cnt int) (
	residue int) {
	var wg sync.WaitGroup
	var idEnqueue, idDequeue int32
	wg.Add(grp)
	q := NewLockFree(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				v := atomic.AddInt32(&idEnqueue, 1)
				ok, _ := q.Enqueue(v)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Enqueue(v)
				}
			}
		}(i)
	}
	wg.Wait()
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < cnt; {
				ok, val, _ := q.Dequeue()
				if !ok {
					fmt.Printf("Dequeue.Fail\n")
					runtime.Gosched()
				} else {
					j++
					idDequeue++
					if idDequeue != val.(int32) {
						t.Logf("Dequeue.Err %d <> %d\n", idDequeue, val)
					}
				}
			}
		}()
	}
	wg.Wait()
	return
}

func TestQueueEnqueueDequeueOrder(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	grp := 1
	cnt := 100

	testQueueEnqueueDequeueOrder(t, grp, cnt)
	t.Logf("Grp: %d, Times: %d", grp, cnt)
}
