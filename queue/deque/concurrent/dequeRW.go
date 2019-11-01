package concurrent

import (
	"goContainer/queue/deque"
	"sync"
)

// DequeRW is implemented on top of deque
//	see it in `https://github.com/Harold2017/goContainer/tree/master/queue/deque`
type DequeRW struct {
	dq   *deque.Deque
	lock sync.RWMutex
}

// NewDequeRW returns a pointer of Deque with capacity >= minCap (8) (defined in deque.deque)
func NewDequeRW(cap int) *DequeRW {
	return &DequeRW{
		dq:   deque.NewDeque(cap),
		lock: sync.RWMutex{},
	}
}

// Cap returns current deque capacity
func (dq *DequeRW) Cap() int {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Cap()
}

// Size returns current deque elements size
func (dq *DequeRW) Size() int {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Size()
}

// Empty returns whether deque is empty
func (dq *DequeRW) Empty() bool {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Empty()
}

// Clear clears deque
//	set all elements inside to nil and rest head / tail pointer, element cnt
func (dq *DequeRW) Clear() {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	dq.dq.Clear()
}

// Values returns elements inside deque's buffer (clockwise)
func (dq *DequeRW) Values() []interface{} {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Values()
}

// String for print
func (dq *DequeRW) String() string {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.String()
}

// PositionsCanPopFront returns the number of positions for pop front in deque
func (dq *DequeRW) PositionsCanPopFront() int {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.PositionsCanPopFront()
}

// PositionsCanPushBack returns the number of available positions for push back in deque
func (dq *DequeRW) PositionsCanPushBack() int {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.PositionsCanPushBack()
}

// IsFull returns true if elements num of deque equal to its capacity
func (dq *DequeRW) IsFull() bool {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.IsFull()
}

// PushBack appends element into deque
func (dq *DequeRW) PushBack(elem interface{}) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	dq.dq.PushBack(elem)
}

// PushFront prepends element into deque
func (dq *DequeRW) PushFront(elem interface{}) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	dq.dq.PushFront(elem)
}

// PopBack returns and delete the last element inside the deque, if deque is empty returns error
func (dq *DequeRW) PopBack() (elem interface{}, err error) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	return dq.dq.PopBack()
}

// PopFront returns and delete the first element inside the deque, if deque is empty returns error
func (dq *DequeRW) PopFront() (elem interface{}, err error) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	return dq.dq.PopFront()
}

// Front returns the first element inside the deque
func (dq *DequeRW) Front() (elem interface{}, err error) {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Front()
}

// Back returns the last element inside the deque
func (dq *DequeRW) Back() (elem interface{}, err error) {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.Back()
}

// At returns element at index idx
func (dq *DequeRW) At(idx int) interface{} {
	dq.lock.RLock()
	defer dq.lock.RUnlock()
	return dq.dq.At(idx)
}
