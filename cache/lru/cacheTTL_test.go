package lru

import (
	"fmt"
	"goContainer/utils"
	"testing"
	"time"
)

func TestNewCacheTTL(t *testing.T) {
	cache, err := NewCacheTTLWithOnEvict(10, 10*time.Second, 0,
		func(key interface{}, value interface{}) {
			fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
		})
	if err != nil {
		t.Fatal(err.Error())
	}
	if cache.Size() != 0 {
		t.Fatal(cache.Size(), cache.lru.size)
	}
	k, v, found := cache.GetLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}
	k, v, found = cache.RemoveLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}

	_, err = NewCacheTTL(-10, 10*time.Second, 0)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestCacheTTL_Add(t *testing.T) {
	cache, err := NewCacheTTL(10, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cache.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	if cache.Size() != 10 {
		t.Fatal(cache.Size(), cache.lru.size)
	}

	for i := range a {
		if !cache.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}

	for i := range a {
		found, evicted := cache.Add(a[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cache.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	cache.Resize(20)
	for i := range a {
		found, evicted := cache.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	cache.Resize(10)
	for i := range b {
		found, evicted := cache.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	k, v, found := cache.GetLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	k, v, found = cache.RemoveLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	time.Sleep(3 * time.Second)
	if cache.Size() != 0 {
		t.Fatal(cache.Size(), cache.lru.size)
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Fatal(cache.Size(), cache.lru.size)
	}
}

func TestCacheTTL_Get(t *testing.T) {
	cache, err := NewCacheTTL(20, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cache.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}
	k, v, found := cache.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cache.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}
	k, v, found = cache.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	for i := range a {
		value, found := cache.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}

	cache.Resize(10)
	for i := range a {
		value, found := cache.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := cache.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := cache.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := cache.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := cache.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}
}

func TestCacheTTL_StopCleanWork(t *testing.T) {
	cache, err := NewCacheTTL(20, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cache.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache.Size(), cache.lru.size)
		}
	}

	cache.StopCleanWork()
	time.Sleep(3 * time.Second)
	if cache.Size() != 10 {
		t.Fatal(cache.Size(), cache.lru.size)
	}

	cache.RestartCleanWork(0)
	time.Sleep(3 * time.Second)
	if cache.Size() != 0 {
		t.Fatal(cache.Size(), cache.lru.size)
	}

	cache.ResetTTL(5 * time.Second)
	if cache.ttl != 5*time.Second {
		t.Fatal(cache.ttl)
	}
}
