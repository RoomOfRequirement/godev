package deque

import (
	"fmt"
)

// minCap represents the minimum capacity of deque
//	deque cap is always power of 2, so that it can utilize bitwise modulus: x % cap = x & (cap - 1), x can be positive or negative
//	e.g. cap = 8, x = -1 -> -1 % 8 == -1 & 7; cap = 8, x = 7 -> 7 % 8 == 7 & 7
//	reference: https://golang.org/src/runtime/slice.go?h=slice#L13
const minCap = 8

// Deque is implemented on top of ring buffer
//	see ring buffer in `https://github.com/Harold2017/ringbuffer/ringbuffer.go`
type Deque struct {
	buf        []interface{} // []interface{} for generics instead of []byte in ringbuffer
	cap, cnt   int
	head, tail int
}

// NewDeque returns a pointer of Deque with capacity >= minCap (8)
func NewDeque(cap int) *Deque {
	if cap <= minCap {
		cap = minCap
	}
	cap = nextPowerOfTwo(cap)
	return &Deque{
		buf:  make([]interface{}, cap),
		cap:  cap,
		cnt:  0,
		head: 0,
		tail: 0,
	}
}

// Cap returns current deque capacity
func (dq *Deque) Cap() int {
	return dq.cap
}

// Size returns current deque elements size
func (dq *Deque) Size() int {
	return dq.cnt
}

// Empty returns whether deque is empty
func (dq *Deque) Empty() bool {
	return dq.cnt == 0
}

// Clear clears deque
//	set all elements inside to nil and rest head / tail pointer, element cnt
func (dq *Deque) Clear() {
	for p := dq.head; p != dq.tail; p = (p + 1) & (dq.cap - 1) {
		dq.buf[p] = nil
	}
	dq.head, dq.tail, dq.cnt = 0, 0, 0
}

// Values returns elements inside deque's buffer (clockwise)
func (dq *Deque) Values() []interface{} {
	if dq.Empty() {
		return nil
	}
	if dq.tail == dq.head {
		buf := make([]interface{}, dq.cap)
		copy(buf, dq.buf)
		return buf
	}
	if dq.tail > dq.head {
		buf := make([]interface{}, dq.tail-dq.head)
		copy(buf, dq.buf[dq.head:dq.tail])
		return buf
	}
	buf := make([]interface{}, dq.cap-dq.head+dq.tail)
	copy(buf, dq.buf[dq.head:dq.cap])
	copy(buf[dq.cap-dq.head:], dq.buf[0:dq.tail])
	return buf
}

// String for print
func (dq *Deque) String() string {
	return fmt.Sprintf("Deque: \n\tCap: %d\n\tPositionsCanPopFront: %d\n\tPositionsCanPushBack: %d\n\tBuffer: %+v\n",
		dq.cap, dq.PositionsCanPopFront(), dq.PositionsCanPushBack(), dq.buf)
}

// PositionsCanPopFront returns the number of positions for pop front in deque
func (dq *Deque) PositionsCanPopFront() int {
	if dq.Empty() {
		return 0
	}
	if dq.tail == dq.head {
		return dq.cap
	}
	if dq.tail > dq.head {
		return dq.tail - dq.head
	}
	return dq.cap - dq.head + dq.tail
}

// PositionsCanPushBack returns the number of available positions for push back in deque
func (dq *Deque) PositionsCanPushBack() int {
	if dq.Empty() {
		return dq.cap
	}
	if dq.tail == dq.head {
		return 0
	}
	if dq.tail < dq.head {
		return dq.head - dq.tail
	}
	return dq.cap - dq.tail + dq.head
}

// IsFull returns true if elements num of deque equal to its capacity
func (dq *Deque) IsFull() bool {
	return dq.cnt == dq.cap
}

func (dq *Deque) resize() {
	// deque size always is double of dq.cnt (half full)
	nCap := dq.cnt << 1
	nBuf := make([]interface{}, nCap)
	if dq.tail > dq.head {
		copy(nBuf, dq.buf[dq.head:dq.tail])
	} else {
		n := copy(nBuf, dq.buf[dq.head:dq.cap])
		copy(nBuf[n:], dq.buf[0:dq.tail])
	}
	dq.cap = nCap
	dq.head = 0
	dq.tail = dq.cnt
	dq.buf = nBuf
}

func (dq *Deque) shrinkCap() {
	// shrink buf cap if buffer is 1 / 4 full
	if dq.cnt<<2 == len(dq.buf) {
		dq.resize()
	}
}

func (dq *Deque) expandCap() {
	// expand buf cap if buffer if full
	if dq.IsFull() {
		dq.resize()
	}
}

// PushBack appends element into deque
func (dq *Deque) PushBack(elem interface{}) {
	dq.expandCap()
	dq.buf[dq.tail] = elem
	dq.tail = (dq.tail + 1) & (dq.cap - 1) // (dq.tail + 1) % dq.cap
	dq.cnt++
}

// PushFront prepends element into deque
func (dq *Deque) PushFront(elem interface{}) {
	dq.expandCap()
	dq.head = (dq.head - 1) & (dq.cap - 1) // (dq.head - 1) % dq.cap
	dq.buf[dq.head] = elem
	dq.cnt++
}

// PopBack returns and delete the last element inside the deque, if deque is empty returns error
func (dq *Deque) PopBack() (elem interface{}, err error) {
	if dq.cnt <= 0 {
		return nil, fmt.Errorf("deque: can NOT PopBack on empty deque")
	}

	dq.tail = (dq.tail - 1) & (dq.cap - 1) // (dq.tail - 1) % dq.cap
	elem = dq.buf[dq.tail]
	dq.buf[dq.tail] = nil
	dq.cnt--

	dq.shrinkCap()
	return
}

// PopFront returns and delete the first element inside the deque, if deque is empty returns error
func (dq *Deque) PopFront() (elem interface{}, err error) {
	if dq.cnt <= 0 {
		return nil, fmt.Errorf("deque: can NOT PopFront on empty deque")
	}

	elem = dq.buf[dq.head]
	dq.buf[dq.head] = nil
	dq.head = (dq.head + 1) & (dq.cap - 1) // (dq.head + 1) % dq.cap
	dq.cnt--

	dq.shrinkCap()
	return
}

// Front returns the first element inside the deque
func (dq *Deque) Front() (elem interface{}, err error) {
	if dq.cnt <= 0 {
		return nil, fmt.Errorf("deque: can NOT get Front element on empty deque")
	}

	elem = dq.buf[dq.head]
	return
}

// Back returns the last element inside the deque
func (dq *Deque) Back() (elem interface{}, err error) {
	if dq.cnt <= 0 {
		return nil, fmt.Errorf("deque: can NOT get Back element on empty deque")
	}

	elem = dq.buf[(dq.tail-1)&(dq.cap-1)]
	return
}

// At returns element at index idx
func (dq *Deque) At(idx int) interface{} {
	if idx < 0 {
		idx = (dq.tail + idx) & (dq.cap - 1)
	} else {
		idx = (dq.head + idx) & (dq.cap - 1)
	}
	return dq.buf[idx]
}

func nextPowerOfTwo(n int) int {
	// https://www.geeksforgeeks.org/smallest-power-of-2-greater-than-or-equal-to-n/
	if n > 0 && n&(n-1) == 0 {
		return n
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16 // 32 bit OS, runtime.GOARCH
	n |= n >> 32 // 64 bit OS
	n++
	return n
}
