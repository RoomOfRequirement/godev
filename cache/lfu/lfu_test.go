package lfu

import (
	"fmt"
	"goContainer/cache"
	"goContainer/utils"
	"testing"
)

func TestNewLFU(t *testing.T) {
	var _ cache.Interface = (*LFU)(nil)

	lfu, err := NewLFU(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	if lfu.Size() != 0 {
		t.Fatal(lfu.Size(), lfu.size)
	}

	k, v, found := lfu.RemoveLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}

	_, err = NewLFU(-10, nil)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestLFU_Add(t *testing.T) {
	lfu, err := NewLFU(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := lfu.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lfu.Size(), lfu.size)
		}
	}

	if lfu.Size() != 10 {
		t.Fatal(lfu.Size(), lfu.size)
	}

	for i := range a {
		if !lfu.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}
	fmt.Println(lfu.freqList.Len())

	b := a[:5]
	//c := a[5:]
	for i := range b {
		found, evicted := lfu.Add(b[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, lfu.Size(), lfu.size)
		}
	}
	if lfu.Size() != 10 {
		t.Fatal(lfu.Size(), lfu.size)
	}
	/*
		fmt.Println(lfu.freqList.Len())
		for i := lfu.freqList.Front(); i != nil; i = i.Next() {
			for it := range i.Value.(*freqItem).items {
				fmt.Printf("%v %v %v\n", i.Value.(*freqItem).freq, it.key, it.value)
			}
			fmt.Println()
		}
	*/

	diff, err := lfu.Resize(20)
	if diff != 10 || err != nil {
		t.Fatal(diff, err)
	}

	diff, err = lfu.Resize(-10)
	if diff != 0 || err == nil {
		t.Fatal(diff, err)
	}

	k, v, found := lfu.RemoveLeastUsed()
	if k == nil || v == nil || !found {
		t.Fatal(k, v, found)
	}

	lfu.Clear()
	if lfu.Size() != 0 {
		t.Fatal(lfu.Size(), lfu.size)
	}
}

func TestTestLRU_Get(t *testing.T) {
	lfu, err := NewLFU(20, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := lfu.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lfu.Size(), lfu.size)
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := lfu.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lfu.Size(), lfu.size)
		}
	}

	for i := range a {
		value, found := lfu.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}

	_, _ = lfu.Resize(10)
	for i := range a {
		value, found := lfu.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := lfu.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := lfu.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := lfu.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := lfu.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}
}
