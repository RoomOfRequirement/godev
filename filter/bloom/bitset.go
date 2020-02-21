package bloomfilter

// BitSet struct
//	byte is the basic unit read by CPU
//	size of golang essential types: https://golang.org/src/go/types/sizes.go
//	uint (32-64)
type BitSet []int64

// NewBitSet creates a new bit set
func NewBitSet(n uint) BitSet {
	// uint just bit shift
	// if int with negative,
	// need consider the 2's complement:
	// (x + ((x >> 31) & ((1 << n) + ^0))) >> n for 32 bit
	// (x + ((x >> 63) & ((1 << n) + ^0))) >> n for 64 bit
	return make(BitSet, (n >> 6) + 1)
}

// SetOne sets idx bit from 0 to 1
//	https://stackoverflow.com/questions/47981/how-do-you-set-clear-and-toggle-a-single-bit
func (bs BitSet) SetOne(idx uint) {
	bitIdx := idx & 63 // bitIdx = idx % 64
	idx = idx >> 6
	bs[idx] |= 1 << bitIdx
}

// SetZero sets idx bit from 1 to 0
func (bs BitSet) SetZero(idx uint) {
	bitIdx := idx & 63
	idx = idx >> 6
	bs[idx] ^= 1 << bitIdx
}

// IsOne returns true if idx-th bit is one
func (bs BitSet) IsOne(idx uint) bool {
	bitIdx := idx & 63
	idx = idx >> 6
	return (bs[idx] | 1<<bitIdx) == bs[idx]
}
