package bufferedchan

import (
	"fmt"
	"goContainer/queue/deque"
	"os"
)

// AutoResize means unlimited buf size
const AutoResize = -1

// BufferedChan struct, just a simple toy, certainly should use make(chan interface, size)
//	a buffered chan = two un-buffered chan + a goroutine + an array
type BufferedChan struct {
	in, out   chan interface{}
	size int
	cap int
	buf       *deque.Deque
}

// New creates a new buffered chan
//	notice: my deque implementation has a min cap = 8
func New(cap int) *BufferedChan {
	if cap < 1 && cap != -1 {
		panic("invalid cap")
	}
	if cap == -1 {
		_, _ = fmt.Fprintln(os.Stderr, "notice: channel buf size is unlimited now!")
	}
	bc := &BufferedChan{
		in:   make(chan interface{}),
		out:  make(chan interface{}),
		size: 0,
		cap:  cap,
		buf:  deque.NewDeque(cap), // >= 8
	}
	go bc.run()
	return bc
}

// In returns a writable chan
func (bc *BufferedChan) In() chan<- interface{} {
	return bc.in
}

// Out returns a readable chan
func (bc *BufferedChan) Out() <-chan interface{} {
	return bc.out
}

// Len returns number of elements in the channel
func (bc *BufferedChan) Len() int {
	return bc.size
}

// Cap returns capacity of the channel
func (bc *BufferedChan) Cap() int {
	return bc.cap
}

// Close closes the channel,
// to keep align with official chan,
// it can read (Out) after close but error when write (In) after close
func (bc *BufferedChan) Close() {
	close(bc.in)
}

func (bc *BufferedChan) run() {
	var in, out, nextIn chan interface{}
	var nextElt interface{}
	// stream replace
	nextIn = bc.in
	in = nextIn

	for in != nil || out != nil {
		select {
			case elt, open := <- in:
				if open {
					bc.buf.PushBack(elt)
					bc.size++
				} else {
					// closed
					in = nil
					nextIn = nil
				}
			case out <- nextElt:
				_, _ = bc.buf.PopFront()
				bc.size--
		}

		if bc.buf.Empty() {
			out = nil
			nextElt = nil
		} else {
			out = bc.out
			nextElt, _ = bc.buf.Front()
		}

		// reach channel capacity
		// if cap is set to AutoResize, it will become an unlimited buf
		if bc.cap != AutoResize && bc.buf.Size() >= bc.cap {
			in = nil
		} else {
			in = nextIn
		}
	}

	close(bc.out)
}
