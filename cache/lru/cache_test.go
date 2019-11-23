package lru

import (
	"fmt"
	"goContainer/utils"
	"math/rand"
	"testing"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCacheWithOnEvict(10, func(key interface{}, value interface{}) {
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

	_, err = NewCacheWithOnEvict(-10, nil)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestCache_Add(t *testing.T) {
	cache, err := NewCache(10)
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

	cache.Clear()
	if cache.Size() != 0 {
		t.Fatal(cache.Size(), cache.lru.size)
	}
}

func TestCache_Get(t *testing.T) {
	cache, err := NewCache(20)
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

func BenchmarkCache(b *testing.B) {
	cache, _ := NewCache(1024 * 8)
	data := make([]int, b.N*2)
	for i := 0; i < b.N*2; i++ {
		data[i] = rand.Intn(1024 * 16)
	}

	b.ResetTimer()

	hit, miss := 0, 0
	for i := 0; i < 2*b.N; i++ {
		if i&1 == 0 { // i % 2
			cache.Add(data[i], data[i])
		} else {
			if _, ok := cache.Get(data[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f\n", hit, miss, float64(hit)/float64(miss))
}
