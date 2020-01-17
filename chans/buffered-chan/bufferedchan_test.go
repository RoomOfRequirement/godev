package bufferedchan

import (
	"godev/chans"
	"godev/utils"
	"testing"
)

func TestBufferedChan(t *testing.T) {
	var _ chans.Interface = (*BufferedChan)(nil)

	bc := New(1000)
	if bc.Cap() != 1000 {
		t.Fatal(bc.Cap())
	}
	// testChanLenCap(t, "New BufferedChan 1000", bc, 1000, 1000)

	bc = New(1000)
	testChan(t, "New BufferedChan 1000", bc)

	bc = New(AutoResize)
	if bc.Cap() != AutoResize {
		t.Fatal(bc.Cap())
	}
	testChan(t, "New BufferedChan unlimited", bc)

	bc = New(100)
	testChan(t, "New BufferedChan 100", bc)

	bc = New(100)
	testChanConcurrent(t, "Concurrent BufferedChan", bc)
}

func testChanLenCap(t *testing.T, name string, ch chans.Interface, expectedLen, expectedCap int) {
	data := make([]int, expectedLen)
	for i := range data {
		data[i] = utils.GenerateRandomInt()
	}

	for _, i := range data {
		ch.In() <- i
	}
	ch.Close()

	if ch.Len() != expectedLen || ch.Cap() != expectedCap {
		t.Fatal("Test ", name, ": expectedLen ", expectedLen, " got ", ch.Len(), ", expectedCap ", expectedCap, " got ", ch.Cap())
	}

	for _, i := range data {
		if num := <-ch.Out(); num.(int) != i {
			t.Fatal("Test ", name, ": expected ", i, " got ", num.(int))
		}
	}
}

func testChan(t *testing.T, name string, ch chans.Interface) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = utils.GenerateRandomInt()
	}

	go func(data []int) {
		for _, i := range data {
			ch.In() <- i
		}
		ch.Close()
	}(data)

	for _, i := range data {
		if num := <-ch.Out(); num.(int) != i {
			t.Fatal("Test ", name, ": expected ", i, " got ", num.(int))
		}
	}
}

func testChanConcurrent(t *testing.T, name string, ch chans.Interface) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = utils.GenerateRandomInt()
	}

	go func(data []int) {
		for _, i := range data {
			ch.In() <- i
		}
		ch.Close()
	}(data)

	go func(data []int) {
		for _, i := range data {
			if num := <-ch.Out(); num.(int) != i {
				t.Fatal("Test ", name, ": expected ", i, " got ", num.(int))
			}
		}
	}(data)

	// check race condition
	go ch.Len()
	go ch.Cap()
}
