package bloomfilter

import (
	"github.com/spaolacci/murmur3"
	"math"
	"strconv"
)

// BloomFilter struct
//	https://en.wikipedia.org/wiki/Bloom_filter
type BloomFilter struct {
	k    uint    // number of hash function
	m    uint    // bit array length
	n    uint    // element number
	p    float64 // false positive rate
	bits BitSet
}

// NewBloomFilter creates a new bloom filter
func NewBloomFilter(n uint, p float64) *BloomFilter {
	k, m := calKM(n, p)
	return &BloomFilter{
		k:    k,
		m:    m,
		n:    n,
		p:    p,
		bits: NewBitSet(m),
	}
}

func calKM(n uint, p float64) (uint, uint) {
	k := -math.Log(p) / math.Log(2.)
	m := k * float64(n) / math.Log(2.)
	return uint(math.Ceil(k)), uint(math.Ceil(m))
}

func hash(data []byte, seed uint) uint {
	m := murmur3.New64WithSeed(uint32(seed))
	_, _ = m.Write(data)
	return uint(m.Sum64())
}

func (bf *BloomFilter) insert(data []byte) {
	l := uint(bf.Size())
	for i := uint(0); i < bf.k; i++ {
		bf.bits.SetOne(hash(data, i) % l)
	}
}

func (bf *BloomFilter) contains(data []byte) bool {
	l := uint(bf.Size())
	for i := uint(0); i < bf.k; i++ {
		if !bf.bits.IsOne(hash(data, i) % l) {
			return false
		}
	}
	return true
}

// InsertKey inserts string key into bloom filter
func (bf *BloomFilter) InsertKey(key string) {
	bf.insert([]byte(key))
}

// ContainsKey returns true if string key inside bloom filter
func (bf *BloomFilter) ContainsKey(key string) bool {
	return bf.contains([]byte(key))
}

// Empty returns true if no key in bloom filter
func (bf *BloomFilter) Empty() bool {
	return len(bf.bits) == 0
}

// Size returns bits number
func (bf *BloomFilter) Size() int {
	return len(bf.bits) << 6
}

// Clear clears bit-set
func (bf *BloomFilter) Clear() {
	bf.bits = nil
}

// Values returns bit-set in string
func (bf *BloomFilter) Values() []interface{} {
	s := make([]interface{}, bf.Size())
	for i := 0; i < bf.Size() >> 6; i++ {
		s[i] = strconv.FormatInt(bf.bits[i], 2)
	}
	return s
}
