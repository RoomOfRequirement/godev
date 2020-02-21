package cuckoofilter

const (
	bucketSize = 4
)

type fingerprint uint16 // 2 bytes

type bucket [bucketSize]fingerprint // 4 fingerprints

var empty fingerprint = 0

func (b *bucket) insert(fp fingerprint) bool {
	for i, p := range b {
		if p == empty {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b *bucket) delete(fp fingerprint) bool {
	for i, p := range b {
		if p == fp {
			b[i] = empty
			return true
		}
	}
	return false
}

func (b *bucket) index(fp fingerprint) int {
	for i, p := range b {
		if p == fp {
			return i
		}
	}
	return -1
}
