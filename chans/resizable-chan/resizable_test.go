package resizablechan

import (
	"goContainer/chans"
	"goContainer/utils"
	"testing"
)

func TestNew(t *testing.T) {
	var _ chans.Interface = (*ResizableChannel)(nil)

	rc := New()
	testChan(t, "New ResizableChannel", rc)

	rc = New()
	err := rc.Resize(AutoResize)
	if err != nil {
		t.Fatal(err)
	}
	testChan(t, "Resize ResizableChannel to unlimited", rc)

	rc = New()
	err = rc.Resize(100)
	if err != nil {
		t.Fatal(err)
	}
	testChan(t, "Resize ResizableChannel to 100", rc)

	rc = New()
	testChanConcurrent(t, "Concurrent ResizableChannel", rc)
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
