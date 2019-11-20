package sort

// Interface is sort interface for all sorting methods, same with golang official sort interface
type Interface interface {
	// Len returns number of elements in the collection
	Len() int
	// Less returns true if element at index i smaller than element at index j
	Less(i, j int) bool
	// Swap swaps element at index i with element at index j
	Swap(i, j int)
}

// IntSlice type for testing sort
type IntSlice []int

// Len returns number of elements in the IntSlice
func (is IntSlice) Len() int {
	return len(is)
}

// Less returns true if element at index i smaller than element at index j
func (is IntSlice) Less(i, j int) bool {
	return is[i] < is[j]
}

// Swap swaps element at index i with element at index j
func (is IntSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func Equal(is1, is2 IntSlice) bool {
	if is1.Len() != is2.Len() {
		return false
	}

	for i := range is1 {
		if is1[i] != is2[i] {
			return false
		}
	}

	return true
}
