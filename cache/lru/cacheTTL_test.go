package lru

import (
	"fmt"
	"goContainer/cache"
	"goContainer/utils"
	"testing"
	"time"
)

func TestNewCacheTTL(t *testing.T) {
	var _ cache.Interface = (*CacheTTL)(nil)

	cacheTTL, err := NewCacheTTLWithOnEvict(10, 10*time.Second, 0,
		func(key interface{}, value interface{}) {
			fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
		})
	if err != nil {
		t.Fatal(err.Error())
	}
	if cacheTTL.Size() != 0 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}
	k, v, found := cacheTTL.GetLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}
	k, v, found = cacheTTL.RemoveLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}

	_, err = NewCacheTTL(-10, 10*time.Second, 0)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestCacheTTL_Add(t *testing.T) {
	cacheTTL, err := NewCacheTTL(10, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cacheTTL.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	if cacheTTL.Size() != 10 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}

	for i := range a {
		if !cacheTTL.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}

	for i := range a {
		found, evicted := cacheTTL.Add(a[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cacheTTL.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	cacheTTL.Resize(20)
	for i := range a {
		found, evicted := cacheTTL.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	cacheTTL.Resize(10)
	for i := range b {
		found, evicted := cacheTTL.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	k, v, found := cacheTTL.GetLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	k, v, found = cacheTTL.RemoveLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	time.Sleep(3 * time.Second)
	if cacheTTL.Size() != 0 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}

	cacheTTL.Clear()
	if cacheTTL.Size() != 0 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}
}

func TestCacheTTL_Get(t *testing.T) {
	cacheTTL, err := NewCacheTTL(20, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cacheTTL.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}
	k, v, found := cacheTTL.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cacheTTL.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}
	k, v, found = cacheTTL.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	for i := range a {
		value, found := cacheTTL.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}

	cacheTTL.Resize(10)
	for i := range a {
		value, found := cacheTTL.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := cacheTTL.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := cacheTTL.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := cacheTTL.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := cacheTTL.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}
}

func TestCacheTTL_StopCleanWork(t *testing.T) {
	cacheTTL, err := NewCacheTTL(20, 1*time.Second, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cacheTTL.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cacheTTL.Size(), cacheTTL.lru.size)
		}
	}

	cacheTTL.StopCleanWork()
	time.Sleep(3 * time.Second)
	if cacheTTL.Size() != 10 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}

	cacheTTL.RestartCleanWork(0)
	time.Sleep(3 * time.Second)
	if cacheTTL.Size() != 0 {
		t.Fatal(cacheTTL.Size(), cacheTTL.lru.size)
	}

	cacheTTL.ResetTTL(5 * time.Second)
	if cacheTTL.ttl != 5*time.Second {
		t.Fatal(cacheTTL.ttl)
	}
}
