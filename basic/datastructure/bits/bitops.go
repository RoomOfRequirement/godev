package bits

// reference: http://graphics.stanford.edu/~seander/bithacks.html

// from `strconv.Iota`
const host32bit = ^uint(0)>>32 == 0
const host64bit = ^uint(0)>>64 == 0
const intSize = 32 << (^uint(0) >> 63)
const shiftToSign = intSize - 1

// const maxInt = int(^uint(0)>>1)

// Sign returns 1 for positive, -1 for negative and 0 for 0
func Sign(x int) int {
	if x == 0 {
		return 0
	}
	return 1 | x>>shiftToSign
}

// IsOppositeSign returns true if x, y have different sign
//	notice: 0 has no sign, which means a 0 in sign flag bit, so it has the same sign with positive number
func IsOppositeSign(x, y int) bool {
	return x^y < 0
}

// Abs returns absolute value of input
func Abs(x int) int {
	sign := x >> shiftToSign
	return x ^ sign - sign // (x + sign) ^ sign
}

// Min return min value of two input ints
//	notice: INT_MIN <= x - y <= INT_MAX
//	if not in this range, better to use simple `>, <` operator for branching
func Min(x, y int) int {
	return y + ((x - y) & ((x - y) >> shiftToSign))
}

// Max return max value of two input ints
//	notice: INT_MIN <= x - y <= INT_MAX
//	if not in this range, better to use simple `>, <` operator for branching
func Max(x, y int) int {
	return x - ((x - y) & ((x - y) >> shiftToSign))
}

// IsPowerOfTwo returns true if x is power of 2
//	notice: 2 ** 0 = 1
func IsPowerOfTwo(x int) bool {
	return x > 0 && x&(x-1) == 0
}

// SetBit sets bit with mask under condition
//	notice: condition = 0, set bit to 0; condition = 1, set bit to 1
//	if (condition) x |= mask; else x &^= mask
func SetBit(x, mask uint, condition uint) uint {
	x = (x &^ mask) | (-condition & mask)
	return x
}

// NegateIf negates x under condition (1 for true, 0 for false)
func NegateIf(x int, condition int) int {
	return (x ^ -condition) + condition
}

// Merge two sets of bits
//	mask: 1 selects bit from y, 0 selects bit from x
func Merge(x, y, mask uint) uint {
	return x ^ ((x ^ y) & mask)
}

// CountBitsOne counts 1s in x's binary representation
func CountBitsOne(x uint) uint {
	var c uint
	/* Native way, for 32 bits uint will loop 32 times
	for c = 0; x != 0; x = x >> 1 {
		c += x & 1
	}
	*/

	// Brian Kernighan's way, loop c times
	for c = 0; x != 0; c++ {
		x &= x - 1 // clear the least significant bit set
	}
	return c
}

// CountBitsOneUntil counts 1s from MSB to position, position starts from LSB
func CountBitsOneUntil(x, position uint) uint {
	return CountBitsOne(x >> position)
}

// ModulusDivide uses bit operator when y is power of 2
func ModulusDivide(x, y int) int {
	if y == 0 {
		panic("zero division")
	}

	if IsPowerOfTwo(y) {
		return x & (y - 1)
	}

	return x % y
}

// CeilToPowerOfTwo returns nearest power of two on ceil
//	notice: x can NOT be 0
//	you can add a branch like `if x ==0 {return 1}` but it is only valid for ceil, it has no floor
func CeilToPowerOfTwo(x int) int {
	if IsPowerOfTwo(x) {
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

// FloorToPowerOfTwo returns nearest power of two on floor
//	notice: x can NOT be 0
func FloorToPowerOfTwo(x int) int {
	if IsPowerOfTwo(x) {
		return x
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x >>= 1
	x++
	return x
}

// CountBitsZeroTailing counts 0s from LSB
/*
0x55555555 = 01010101 01010101 01010101 01010101
0x33333333 = 00110011 00110011 00110011 00110011
0x0F0F0F0F = 00001111 00001111 00001111 00001111
0x00FF00FF = 00000000 11111111 00000000 11111111
0x0000FFFF = 00000000 00000000 11111111 11111111
*/
func CountBitsZeroTailing(x uint) int {
	if host64bit {
		c := 64
		x &= -x
		if x != 0 {
			c--
		}
		if x&0x00000000FFFFFFFF != 0 {
			c -= 32
		}
		if x&0x0000FFFF0000FFFF != 0 {
			c -= 16
		}
		if x&0x00FF00FF00FF00FF != 0 {
			c -= 8
		}
		if x&0x0F0F0F0F0F0F0F0F != 0 {
			c -= 4
		}
		if x&0x3333333333333333 != 0 {
			c -= 2
		}
		if x&0x5555555555555555 != 0 {
			c -= 1
		}
		return c
	}
	if host32bit {
		c := 32
		x &= -x
		if x != 0 {
			c--
		}
		if x&0x0000FFFF != 0 {
			c -= 16
		}
		if x&0x00FF00FF != 0 {
			c -= 8
		}
		if x&0x0F0F0F0F != 0 {
			c -= 4
		}
		if x&0x33333333 != 0 {
			c -= 2
		}
		if x&0x55555555 != 0 {
			c -= 1
		}
		return c
	}
	panic("arch not support")
}

// IsOdd ...
func IsOdd(x int) bool {
	// or `x % 2 != 0`, they almost have the same performance,
	// see benchmark in `bitopes_test.go`
	// see compiler explorer: https://gcc.godbolt.org/#g:!((g:!((g:!((h:codeEditor,i:(j:1,options:(compileOnChange:'0'),source:'int+isOdd_mod(unsigned+x)+%7B%0A++++return+(x+%25+2)%3B%0A%7D%0A%0Aint+isOdd_and(unsigned+x)+%7B%0A++++return+(x+%26+1)%3B%0A%7D%0A%0Aint+isOdd_or(unsigned+x)+%7B%0A++++return+(0xFFFFFFFF+%3D%3D+(x+%7C+0xFFFFFFFE))%3B%0A%7D+++'),l:'5',n:'1',o:'C%2B%2B+source+%231',t:'0')),k:33.333333333333336,l:'4',n:'0',o:'',s:0,t:'0'),(g:!((h:compiler,i:(compiler:clang390,filters:(b:'0',commentOnly:'0',directives:'0',intel:'0'),options:'-O3'),l:'5',n:'0',o:'%231+with+x86-64+clang+3.9.0',t:'0')),k:33.333333333333336,l:'4',n:'0',o:'',s:0,t:'0'),(g:!((h:compiler,i:(compiler:g62,filters:(b:'0',commentOnly:'0',directives:'0',intel:'0'),options:'-O3'),l:'5',n:'0',o:'%231+with+x86-64+gcc+6.2',t:'0')),k:33.33333333333333,l:'4',n:'0',o:'',s:0,t:'0')),l:'2',n:'0',o:'',t:'0')),version:4
	return (x & 1) == 1
}

// IsEven ...
func IsEven(x int) bool {
	return (x & 1) == 0
}
