package byteslicepool

import (
	"godev/basic/datastructure/bits"
	"sync"
	"sync/atomic"
)

type pool struct {
	buckets *buckets // cap -> bucket
	metrics Metrics
}

// Get gets one byte slice with `size` or create a new one
func (p *pool) Get(size int) []byte {
	atomic.AddUint64(&p.metrics.gets, 1)
	bCap := bits.CeilToPowerOfTwo(size) // power of 2, aligned with gc
	bs := p.buckets.Get(bCap, size)
	if bs == nil {
		return make([]byte, size, bCap) // make a new one
	}
	atomic.AddUint64(&p.metrics.hits, 1) // hit
	atomic.AddInt64(&p.metrics.size, int64(-bCap))
	return bs
}

func (p *pool) Put(byteSlice []byte) {
	bCap := bits.CeilToPowerOfTwo(len(byteSlice))
	atomic.AddUint64(&p.metrics.puts, 1)
	atomic.AddInt64(&p.metrics.size, int64(bCap))
	p.buckets.Put(bCap, byteSlice)
}

func (p *pool) Metrics() Metrics {
	return Metrics{
		gets: atomic.LoadUint64(&p.metrics.gets),
		puts: atomic.LoadUint64(&p.metrics.puts),
		hits: atomic.LoadUint64(&p.metrics.hits),
		size: atomic.LoadInt64(&p.metrics.size),
	}
}

type buckets struct {
	p map[int]*sync.Pool
}

func (b *buckets) Get(cap, size int) []byte {
	if bucket := b.p[cap]; bucket != nil {
		return bucket.Get().([]byte)[:size] // copy
	}
	return nil
}

func (b *buckets) Put(cap int, byteSlice []byte) {
	// only one bucket with certain cap
	if bucket := b.p[cap]; bucket != nil {
		bucket.Put(byteSlice)
	} else {
		b.p[cap] = new(sync.Pool)
		b.p[cap].Put(byteSlice)
	}
}
