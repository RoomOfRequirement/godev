package resizablechan

import (
	"errors"
	"fmt"
	"goContainer/basic/datastructure/queue/deque"
	"os"
)

// AutoResize means unlimited buf size
const AutoResize = -1

// ResizableChannel struct
//	reference: https://github.com/eapache/channels
type ResizableChannel struct {
	in, out          chan interface{}
	len, cap, resize chan int
	bufSize          int
	buf              *deque.Deque
}

// New creates a new resizable channel
func New() *ResizableChannel {
	rc := &ResizableChannel{
		in:      make(chan interface{}),
		out:     make(chan interface{}),
		len:     make(chan int),
		cap:     make(chan int),
		resize:  make(chan int),
		bufSize: 1,
		buf:     deque.NewDeque(1), // deque minCap = 8
	}
	go rc.run()
	return rc
}

// In returns a writable chan
func (rc *ResizableChannel) In() chan<- interface{} {
	return rc.in
}

// Out returns a readable chan
func (rc *ResizableChannel) Out() <-chan interface{} {
	return rc.out
}

// Len returns number of elements in the channel
func (rc *ResizableChannel) Len() int {
	return <-rc.len
}

// Cap returns capacity of the channel
func (rc *ResizableChannel) Cap() int {
	cCap, open := <-rc.cap
	if open {
		return cCap
	}
	return rc.bufSize
}

// Resize resize channel
// size -> rc.resize (-> rc.bufSize) -> rc.len -> rc.cap
func (rc *ResizableChannel) Resize(size int) error {
	if size <= 0 && size != AutoResize {
		return errors.New("invalid channel size")
	}
	if size == AutoResize {
		_, _ = fmt.Fprintln(os.Stderr, "notice: channel buf size is unlimited now!")
	}
	rc.resize <- size
	return nil
}

// Close closes the channel,
// to keep align with official chan,
// it can read (Out) after close but error when write (In) after close
func (rc *ResizableChannel) Close() {
	close(rc.in)
}

// run resizable channel in a background goroutine
// get in element and put it into buf
// get out element from buf and put it into out
// update len, cap according to buf size
// resize channel according to resize chan
func (rc *ResizableChannel) run() {
	var in, out, nextIn chan interface{}
	var nextElt interface{}
	// stream replace
	nextIn = rc.in
	in = nextIn

	for in != nil || out != nil {
		select {
		case elt, open := <-in:
			if open {
				rc.buf.PushBack(elt)
			} else {
				// closed
				in = nil
				nextIn = nil
			}
		case out <- nextElt:
			// two steps for out: 2. pop
			_, _ = rc.buf.PopFront()
		case rc.bufSize = <-rc.resize:
		case rc.len <- rc.buf.Size():
		case rc.cap <- rc.bufSize:
		}

		// no elt in buf -> nothing to out
		if rc.buf.Empty() {
			out = nil
			nextElt = nil
		} else {
			out = rc.out
			// to avoid unnecessary resize check of buf
			// two steps for out: 1. pre-read
			nextElt, _ = rc.buf.Front()
		}

		// reach channel capacity
		// if cap is set to AutoResize, it will become an unlimited buf
		if rc.bufSize != AutoResize && rc.buf.Size() >= rc.bufSize {
			in = nil
		} else {
			in = nextIn
		}
	}

	close(rc.out)
	close(rc.resize)
	close(rc.len)
	close(rc.cap)
}
