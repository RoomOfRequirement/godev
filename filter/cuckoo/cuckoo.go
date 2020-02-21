package cuckoofilter

import "math/rand"

// MaxNumKicks ... max trials for reinsertion
const MaxNumKicks = 500

// CuckooFilter ...
//	based on https://github.com/seiflotfy/cuckoofilter
/*
 * Cuckoo Filter: Practically Better Than Bloom by Bin Fan, David G. Andersen, Michael Kaminsky, Michael D. Mitzenmacher (2014)
 * http://www.eecs.harvard.edu/~michaelm/postscripts/cuckoo-conext2014.pdf
 * (2,4)-cuckoo filter: 2 hash, 4 bucket size
 */
type CuckooFilter struct {
	buckets []bucket // buckets
	n       int      // element number
}

// NewCuckooFilter ...
func NewCuckooFilter(capacity uint) *CuckooFilter {
	// must be power of 2
	capacity = ceilToPowerOfTwo(capacity) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	// initialize
	buckets := make([]bucket, capacity)
	for i := range buckets {
		buckets[i] = bucket{}
	}
	return &CuckooFilter{
		buckets: buckets,
		n:       0,
	}
}

func (cf *CuckooFilter) insert(data []byte) bool {
	fp := genFingerprint(data)
	idx1 := genIdx(data, uint(len(cf.buckets)))
	idx2 := genAltIdx(fp, idx1, uint(len(cf.buckets)))
	if cf.buckets[idx1].insert(fp) || cf.buckets[idx2].insert(fp) {
		cf.n++
		return true
	}
	return cf.reinsert(fp, idx2)
}

func (cf *CuckooFilter) reinsert(fp fingerprint, idx uint) bool {
	for i := 0; i < MaxNumKicks; i++ {
		r := rand.Intn(bucketSize)
		oldFp := fp
		fp = cf.buckets[idx][r]
		cf.buckets[idx][r] = oldFp
		idx = genAltIdx(fp, idx, uint(len(cf.buckets)))
		if cf.buckets[idx].insert(fp) {
			cf.n++
			return true
		}
	}
	return false
}

func (cf *CuckooFilter) delete(data []byte) bool {
	fp := genFingerprint(data)
	idx1 := genIdx(data, uint(len(cf.buckets)))
	idx2 := genAltIdx(fp, idx1, uint(len(cf.buckets)))
	if cf.buckets[idx1].delete(fp) || cf.buckets[idx2].delete(fp) {
		cf.n--
		return true
	}
	return false
}

func (cf *CuckooFilter) lookup(data []byte) bool {
	fp := genFingerprint(data)
	idx1 := genIdx(data, uint(len(cf.buckets)))
	idx2 := genAltIdx(fp, idx1, uint(len(cf.buckets)))
	return cf.buckets[idx1].index(fp) != -1 || cf.buckets[idx2].index(fp) != -1
}

// InsertKey ...
func (cf *CuckooFilter) InsertKey(key string) bool {
	return cf.insert([]byte(key))
}

// ContainsKey ...
func (cf *CuckooFilter) ContainsKey(key string) bool {
	return cf.lookup([]byte(key))
}

// DeleteKey ...
func (cf *CuckooFilter) DeleteKey(key string) bool {
	return cf.delete([]byte(key))
}

// Empty ...
func (cf *CuckooFilter) Empty() bool {
	return cf.n == 0
}

// Size ...
func (cf *CuckooFilter) Size() int {
	return cf.n
}

// Clear ...
func (cf *CuckooFilter) Clear() {
	cf.buckets = nil
	cf.n = 0
}

// Values ...
func (cf *CuckooFilter) Values() []interface{} {
	bytes := make([]interface{}, len(cf.buckets)*bucketSize)
	for i, b := range cf.buckets {
		for j, f := range b {
			index := (i * len(b)) + j
			bytes[index] = f
		}
	}
	return bytes
}
