package cache

import (
	"fmt"
	"goContainer/cache"
	"goContainer/utils"
	"math/rand"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	var _ cache.Interface = (*Cache)(nil)

	types := []Type{LRU, LFU, ARC}

	for _, ty := range types {
		cache2, err := NewCacheWithOnEvict(10, ty, func(key interface{}, value interface{}) {
			fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
		})
		if err != nil {
			t.Fatal(err.Error())
		}
		if cache2.Size() != 0 {
			t.Fatal(cache2.Size(), cache2.cache)
		}

		k, v, found := cache2.RemoveLeastUsed()
		if k != nil || v != nil || found {
			t.Fatal(k, v, found)
		}

		_, err = NewCacheWithOnEvict(-10, ty, nil)
		if err == nil || err.Error() != "invalid cache size" {
			t.Fatal(err)
		}
	}
}

// Add function of LRU, LFU, ARC is not consistent!
func TestCache_Add(t *testing.T) {
	cache2, err := NewCache(10, LRU)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cache2.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	if cache2.Size() != 10 {
		t.Fatal(cache2.Size())
	}

	for i := range a {
		if !cache2.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}

	for i := range a {
		found, evicted := cache2.Add(a[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cache2.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	cache2.Resize(20)
	for i := range a {
		found, evicted := cache2.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	cache2.Resize(10)
	for i := range b {
		found, evicted := cache2.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	k, v, found := cache2.RemoveLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	cache2.Clear()
	if cache2.Size() != 0 {
		t.Fatal(cache2.Size())
	}
}

// Get function (after Add) of LRU, LFU, ARC is not consistent!
func TestCache_Get(t *testing.T) {
	cache2, err := NewCache(20, LRU)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := cache2.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := cache2.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, cache2.Size())
		}
	}
	for i := range a {
		value, found := cache2.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}

	cache2.Resize(10)
	for i := range a {
		value, found := cache2.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := cache2.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := cache2.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := cache2.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := cache2.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}
}

// Benchmark reference: https://github.com/hashicorp/golang-lru/blob/master/arc_test.go
func init() {
	rand.Seed(time.Now().Unix())
}

func BenchmarkCacheRandom(b *testing.B) {
	types := map[string]Type{"LRU": LRU, "LFU": LFU, "ARC": ARC}

	for ky, ty := range types {
		b.Run("cache type: "+ky, func(b *testing.B) {
			cache2, _ := NewCache(1024*8, ty)
			data := make([]int64, b.N*2)
			for i := 0; i < b.N*2; i++ {
				data[i] = rand.Int63() % 32768
			}

			b.ResetTimer()

			hit, miss := 0, 0
			for i := 0; i < 2*b.N; i++ {
				if i&1 == 0 {
					cache2.Add(data[i], data[i])
				} else {
					if _, ok := cache2.Get(data[i]); ok {
						hit++
					} else {
						miss++
					}
				}
			}
			b.Logf("hit: %d miss: %d ratio: %f\n", hit, miss, float64(hit)/float64(miss))
		})
	}
}

/*
BenchmarkCacheRandom/cache_type:_LRU-8         	 3000000	       443 ns/op
--- BENCH: BenchmarkCacheRandom/cache_type:_LRU-8
    cache_test.go:204: hit: 0 miss: 1 ratio: 0.000000
    cache_test.go:204: hit: 1 miss: 99 ratio: 0.010101
    cache_test.go:204: hit: 1450 miss: 8550 ratio: 0.169591
    cache_test.go:204: hit: 248950 miss: 751050 ratio: 0.331469
    cache_test.go:204: hit: 750465 miss: 2249535 ratio: 0.333609
BenchmarkCacheRandom/cache_type:_LFU-8         	 2000000	       642 ns/op
--- BENCH: BenchmarkCacheRandom/cache_type:_LFU-8
    cache_test.go:204: hit: 0 miss: 1 ratio: 0.000000
    cache_test.go:204: hit: 0 miss: 100 ratio: 0.000000
    cache_test.go:204: hit: 1403 miss: 8597 ratio: 0.163196
    cache_test.go:204: hit: 249045 miss: 750955 ratio: 0.331638
    cache_test.go:204: hit: 497706 miss: 1502294 ratio: 0.331297
BenchmarkCacheRandom/cache_type:_ARC-8         	 2000000	       722 ns/op
--- BENCH: BenchmarkCacheRandom/cache_type:_ARC-8
    cache_test.go:204: hit: 0 miss: 1 ratio: 0.000000
    cache_test.go:204: hit: 0 miss: 100 ratio: 0.000000
    cache_test.go:204: hit: 1320 miss: 8680 ratio: 0.152074
    cache_test.go:204: hit: 248621 miss: 751379 ratio: 0.330886
    cache_test.go:204: hit: 498642 miss: 1501358 ratio: 0.332127
*/

func BenchmarkCacheFrequent(b *testing.B) {
	types := map[string]Type{"LRU": LRU, "LFU": LFU, "ARC": ARC}

	for ky, ty := range types {
		b.Run("cache type: "+ky, func(b *testing.B) {
			cache2, _ := NewCache(1024*8, ty)
			data := make([]int64, b.N*2)
			for i := 0; i < b.N*2; i++ {
				if i&1 == 0 {
					data[i] = rand.Int63() % 16384
				} else {
					data[i] = rand.Int63() % 32768
				}
			}

			b.ResetTimer()

			hit, miss := 0, 0
			for i := 0; i < b.N; i++ {
				cache2.Add(data[i], data[i])
			}
			for i := 0; i < b.N; i++ {
				if _, ok := cache2.Get(data[i]); ok {
					hit++
				} else {
					miss++
				}
			}
			b.Logf("hit: %d miss: %d ratio: %f\n", hit, miss, float64(hit)/float64(miss))
		})
	}
}

/*
BenchmarkCacheFrequent/cache_type:_LRU-8         	 3000000	       390 ns/op
--- BENCH: BenchmarkCacheFrequent/cache_type:_LRU-8
    cache_test.go:262: hit: 1 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 100 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 9851 miss: 149 ratio: 66.114094
    cache_test.go:262: hit: 314497 miss: 685503 ratio: 0.458783
    cache_test.go:262: hit: 929352 miss: 2070648 ratio: 0.448822
BenchmarkCacheFrequent/cache_type:_LFU-8         	 2000000	       691 ns/op
--- BENCH: BenchmarkCacheFrequent/cache_type:_LFU-8
    cache_test.go:262: hit: 1 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 100 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 9851 miss: 149 ratio: 66.114094
    cache_test.go:262: hit: 347793 miss: 652207 ratio: 0.533256
    cache_test.go:262: hit: 687815 miss: 1312185 ratio: 0.524175
BenchmarkCacheFrequent/cache_type:_ARC-8         	 2000000	       749 ns/op
--- BENCH: BenchmarkCacheFrequent/cache_type:_ARC-8
    cache_test.go:262: hit: 1 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 100 miss: 0 ratio: +Inf
    cache_test.go:262: hit: 9849 miss: 151 ratio: 65.225166
    cache_test.go:262: hit: 351704 miss: 648296 ratio: 0.542505
    cache_test.go:262: hit: 692739 miss: 1307261 ratio: 0.529916
*/
