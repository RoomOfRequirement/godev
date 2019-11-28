package chans

// Interface means channel interface
//	https://golang.org/src/runtime/chan.go
type Interface interface {
	// In means writable
	In() chan<- interface{}
	// Out means readable
	Out() <-chan interface{}
	// Len returns number of elements in the channel
	Len() int
	// Cap returns capacity of the channel
	Cap() int
	// Close closes the channel,
	// to keep align with official chan,
	// it can read after close but error when write after close
	Close()
}
