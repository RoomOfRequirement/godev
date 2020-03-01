package kmp

// my naive implementation
type kmp struct {
	states  [][]int32
	pattern string
}

func newKMP(pattern string) *kmp {
	m := int32(len(pattern))
	k := &kmp{
		states:  nil,
		pattern: pattern,
	}
	states := make([][]int32, m)
	for i := range states {
		states[i] = make([]int32, 256) // 256 for ASCII
	}
	states[0][rune(pattern[0])] = 1
	// x: shadow state
	// c: rune at pattern[j]
	var x, j, c int32
	for j = 1; j < m; j++ {
		for c = 0; c < 256; c++ {
			if rune(pattern[j]) == c {
				states[j][c] = j + 1 // next moves forward
			} else {
				states[j][c] = states[x][c] // next moves backward, use shadow state x to reduce backwards
			}
			x = states[x][rune(pattern[j])]
		}
	}
	k.states = states
	return k
}

func (k *kmp) search(txt string) int {
	m, n := len(k.pattern), len(txt)
	var j int32 = 0
	for i := 0; i < n; i++ {
		j = k.states[j][rune(txt[i])]
		if j == int32(m) {
			return i - m + 1
		}
	}
	return -1
}

// KMP ...
//	https://en.wikipedia.org/wiki/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm
type KMP struct {
	next    []int
	pattern string
}

// New ...
func New(pattern string) *KMP {
	m := len(pattern)
	k := &KMP{
		next:    nil,
		pattern: pattern,
	}
	next := make([]int, m+1)
	i, j := 0, -1
	next[0] = -1
	for i < m-1 {
		for j > -1 && pattern[i] != pattern[j] {
			j = next[j]
		}
		i++
		j++
		if pattern[i] == pattern[j] {
			next[i] = next[j]
		} else {
			next[i] = j
		}
	}
	k.next = next
	return k
}

// Search ...
//	all match idx array
func (k *KMP) Search(txt string) []int {
	i, j := 0, 0
	m, n := len(k.pattern), len(txt)
	x, y := []byte(k.pattern), []byte(txt)
	var ret []int

	if m == 0 || n == 0 || n < m {
		return ret
	}

	for j < n {
		for i > -1 && x[i] != y[j] {
			i = k.next[i]
		}
		i++
		j++
		// found
		if i >= m {
			ret = append(ret, j-i)
			i = k.next[i]
		}
	}
	return ret
}

// SearchFirst ...
//	first position (first match)
func (k *KMP) SearchFirst(txt string) int {
	ret := k.Search(txt)
	if len(ret) > 0 {
		return ret[0]
	}
	return -1
}

// SearchLast ...
//	last position (last match)
func (k *KMP) SearchLast(txt string) int {
	ret := k.Search(txt)
	if len(ret) > 0 {
		return ret[len(ret)-1]
	}
	return -1
}
