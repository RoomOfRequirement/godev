package cuckoofilter

import (
	"bytes"
	"encoding/binary"
	"github.com/spaolacci/murmur3"
)

// HashFunc ... uint64 -> 8 bytes
type HashFunc func([]byte) uint64

// DefaultHash func used in fingerprint generation and idx generation
var DefaultHash HashFunc = murmur3.Sum64

func ceilToPowerOfTwo(x uint) uint {
	if x > 0 && x&(x-1) == 0 {
		return x
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}

// fingerprint of data is a reduced bit string of an input string
func genFingerprint(data []byte) fingerprint {
	// 2 bytes MSB for fingerprint, uint16
	return fingerprint(DefaultHash(data) >> (64 - 16))
}

func genIdx(data []byte, numBuckets uint) uint {
	return uint(DefaultHash(data)) % numBuckets
}

func genAltIdx(fp fingerprint, idx1 uint, numBuckets uint) uint {
	// xor
	return (idx1 ^ genIdx(fingerprint2bytes(fp), numBuckets)) % numBuckets
}

func fingerprint2bytes(fp fingerprint) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, fp)
	return buff.Bytes()
}
