package queue

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// minCap of queue
const minCap = 8

// BLK storage block in queue
type BLK struct {
	hN, tN uint32
	val    interface{}
}

// LockFree queue based on ring-buffer
//	notice: FIFO
type LockFree struct {
	cap        uint32
	head, tail uint32
	buf        []BLK
}

// NewLockFree ...
func NewLockFree(cap uint32) *LockFree {
	if cap <= minCap {
		cap = minCap
	}
	cap = nextPowerOfTwo(cap)
	lf := &LockFree{
		cap:  cap,
		head: 0,
		tail: 0,
		buf:  make([]BLK, cap),
	}
	for i := range lf.buf {
		lf.buf[i].hN = uint32(i)
		lf.buf[i].tN = uint32(i)
	}
	lf.buf[0].hN = cap
	lf.buf[0].tN = cap
	return lf
}

// Cap returns current queue capacity
func (lf *LockFree) Cap() uint32 {
	return lf.cap
}

// Size returns current queue elements size
//	notice: max size = cap - 1
func (lf *LockFree) Size() uint32 {
	var head, tail, size uint32
	head = atomic.LoadUint32(&lf.head)
	tail = atomic.LoadUint32(&lf.tail)
	// tail < head when tail exceed uint32 boundary
	if tail >= head {
		size = tail - head
	} else {
		size = lf.cap + tail - head - 1
	}
	return size
}

// Empty returns whether queue is empty
func (lf *LockFree) Empty() bool {
	return lf.Size() == 0
}

// String for print
func (lf *LockFree) String() string {
	return fmt.Sprintf("Queue: \n\tCap: %d\n\tSize: %d\n\tBuffer: %+v\n",
		lf.cap, lf.Size(), lf.buf)
}

// Enqueue ...
func (lf *LockFree) Enqueue(val interface{}) (success bool, size uint32) {
	var head, tail, newTail, cnt uint32
	var blk *BLK
	capMod := lf.cap - 1
	head = atomic.LoadUint32(&lf.head)
	tail = atomic.LoadUint32(&lf.tail)
	if tail >= head {
		cnt = tail - head
	} else {
		cnt = capMod + tail - head
	}
	if cnt >= capMod {
		runtime.Gosched()
		return false, cnt
	}

	newTail = tail + 1
	if !atomic.CompareAndSwapUint32(&lf.tail, tail, newTail) {
		runtime.Gosched()
		return false, cnt
	}

	blk = &lf.buf[newTail&capMod]
	for {
		hN := atomic.LoadUint32(&blk.hN)
		tN := atomic.LoadUint32(&blk.tN)
		if newTail == tN && hN == tN {
			blk.val = val
			atomic.AddUint32(&blk.tN, capMod+1)
			return true, cnt + 1
		}
		runtime.Gosched()
	}
}

// Dequeue ...
func (lf *LockFree) Dequeue() (success bool, val interface{}, size uint32) {
	var head, tail, newHead, cnt uint32
	var blk *BLK
	capMod := lf.cap - 1
	head = atomic.LoadUint32(&lf.head)
	tail = atomic.LoadUint32(&lf.tail)
	if tail >= head {
		cnt = tail - head
	} else {
		cnt = capMod + tail - head
	}
	if cnt < 1 {
		runtime.Gosched()
		return false, nil, cnt
	}

	newHead = head + 1
	if !atomic.CompareAndSwapUint32(&lf.head, head, newHead) {
		runtime.Gosched()
		return false, nil, cnt
	}

	blk = &lf.buf[newHead&capMod]
	for {
		hN := atomic.LoadUint32(&blk.hN)
		tN := atomic.LoadUint32(&blk.tN)
		if newHead == hN && hN == tN-capMod-1 {
			val = blk.val
			blk.val = nil
			atomic.AddUint32(&blk.hN, capMod+1)
			return true, val, cnt - 1
		}
		runtime.Gosched()
	}
}

/*
func (lf *LockFree) EnqueueN(vals []interface{}) (n, size uint32) {
	var head, tail, newTail, cnt uint32
	var blk *BLK
	capMod := lf.cap - 1
	head = atomic.LoadUint32(&lf.head)
	tail = atomic.LoadUint32(&lf.tail)
	if tail >= head {
		cnt = tail - head
	} else {
		cnt = capMod + tail - head
	}
	if cnt >= capMod {
		runtime.Gosched()
		return 0, cnt
	}

	if n, m := capMod + 1 - cnt, uint32(len(vals)); n >= m {
		n = m
	}

	newTail = tail + n
	if !atomic.CompareAndSwapUint32(&lf.tail, tail, newTail) {
		runtime.Gosched()
		return 0, cnt
	}

	for t, i := tail+1, uint32(0); i < n; t, i = t + 1, i + 1 {
		blk = &lf.buf[t & capMod]
		for {
			hN := atomic.LoadUint32(&blk.hN)
			tN := atomic.LoadUint32(&blk.tN)
			if t == tN && hN == tN {
				blk.val = vals[i]
				atomic.AddUint32(&blk.tN, capMod + 1)
				break
			}
			runtime.Gosched()
		}
	}
	return n, cnt + n
}

func (lf *LockFree) DequeueN(n uint32) (m uint32, vp *[]interface{}, size uint32) {
	var head, tail, newHead, cnt uint32
	var blk *BLK
	capMod := lf.cap - 1
	head = atomic.LoadUint32(&lf.head)
	tail = atomic.LoadUint32(&lf.tail)
	if tail >= head {
		cnt = tail - head
	} else {
		cnt = capMod + tail - head
	}
	if cnt < 1 {
		runtime.Gosched()
		return 0, nil, cnt
	}

	if cnt >= n {
		m = n
	} else {
		m = cnt
	}

	newHead = head + m
	if !atomic.CompareAndSwapUint32(&lf.head, head, newHead) {
		runtime.Gosched()
		return 0, nil, cnt
	}

	vals := make([]interface{}, m)
	for h, i := head+1, uint32(0); i < m; h, i = h + 1, i + 1 {
		blk = &lf.buf[h&capMod]
		for {
			hN := atomic.LoadUint32(&blk.hN)
			tN := atomic.LoadUint32(&blk.tN)
			if h == hN && hN == tN - capMod - 1 {
				vals[i] = blk.val
				blk.val = nil
				atomic.AddUint32(&blk.hN, capMod + 1)
				break
			}
			runtime.Gosched()
		}
	}

	return m, &vals, cnt - m
}
*/

func nextPowerOfTwo(n uint32) uint32 {
	if n > 0 && n&(n-1) == 0 {
		return n
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
