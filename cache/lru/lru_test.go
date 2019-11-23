package lru

import (
	"fmt"
	"goContainer/utils"
	"testing"
)

func TestNewLRU(t *testing.T) {
	lru, err := NewLRU(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if lru.Size() != 0 {
		t.Fatal(lru.Size(), lru.size)
	}
	k, v, found := lru.GetLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}
	k, v, found = lru.RemoveLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}

	_, err = NewLRU(-10, nil)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestLRU_Add(t *testing.T) {
	lru, err := NewLRU(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := lru.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}

	if lru.Size() != 10 {
		t.Fatal(lru.Size(), lru.size)
	}

	for i := range a {
		if !lru.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}

	for i := range a {
		found, evicted := lru.Add(a[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := lru.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}

	lru.Resize(20)
	for i := range a {
		found, evicted := lru.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}

	lru.Resize(10)
	for i := range b {
		found, evicted := lru.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}

	k, v, found := lru.GetLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	k, v, found = lru.RemoveLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	lru.Clear()
	if lru.Size() != 0 {
		t.Fatal(lru.Size(), lru.size)
	}
}

func TestLRU_Get(t *testing.T) {
	lru, err := NewLRU(20, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := lru.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}
	k, v, found := lru.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := lru.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, lru.Size(), lru.size)
		}
	}
	k, v, found = lru.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	for i := range a {
		value, found := lru.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}

	lru.Resize(10)
	for i := range a {
		value, found := lru.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := lru.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := lru.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := lru.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := lru.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}
}
