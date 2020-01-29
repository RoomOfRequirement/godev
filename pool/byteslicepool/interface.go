package byteslicepool

import (
	"fmt"
	"sync"
)

// Pool interface
type Pool interface {
	Get(size int) []byte
	Put(byteSlice []byte)
	Metrics() Metrics
}

// Metrics for investigation
type Metrics struct {
	gets uint64
	puts uint64
	hits uint64
	size int64
}

// String for print
func (m Metrics) String() string {
	return fmt.Sprintf("Byte Slice Pool Metrics: \n\tGets: %d\n\tPuts: %d\n\tHits: %d\n\tSize: %d\n", m.gets, m.puts, m.hits, m.size)
}

// New creates a new byte slice pool
func New() Pool {
	return &pool{
		buckets: &buckets{
			p: make(map[int]*sync.Pool),
		},
		metrics: Metrics{
			gets: 0,
			puts: 0,
			hits: 0,
			size: 0,
		},
	}
}
