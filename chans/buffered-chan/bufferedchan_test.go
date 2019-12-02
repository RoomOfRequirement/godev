package bufferedchan

import (
	"goContainer/chans"
	"goContainer/utils"
	"testing"
)

func TestBufferedChan(t *testing.T) {
	var _ chans.Interface = (*BufferedChan)(nil)

	bc := New(1000)
	testChan(t, "New BufferedChan", bc)

	bc = New(AutoResize)
	testChan(t, "New BufferedChan unlimited", bc)

	bc = New(100)
	testChan(t, "New BufferedChan 100", bc)

	bc = New(100)
	testChanConcurrent(t, "Concurrent BufferedChan", bc)
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

	go ch.Len()
	go ch.Cap()
}
