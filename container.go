package container

import (
	"sort"
)

// Container interface
type Container interface {
	Size() int
	Empty() bool
	Clear()
	Values() []interface{}
	// Sort()
	// SortStable()
}

// Comparator for sorting
//	 1 if a > b
//	 0 if a == b
//	-1 if a < b
//	TODO: i think i have to make many type assertions... any better way?
type Comparator func(a, b interface{}) int

// IntComparator compares two ints
func IntComparator(a, b interface{}) int {
	A := a.(int)
	B := b.(int)
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

// Float64Comparator compares two float64s
func Float64Comparator(a, b interface{}) int {
	A := a.(float64)
	B := b.(float64)
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

// sortable obeys sort.Interface, which requires Len() / Swap() / Less()
type sortable struct {
	values     []interface{}
	comparator Comparator
}

func (st sortable) Len() int {
	return len(st.values)
}

func (st sortable) Swap(i, j int) {
	st.values[i], st.values[j] = st.values[j], st.values[i]
}

func (st sortable) Less(i, j int) bool {
	return st.comparator(st.values[i], st.values[j]) < 0
}

// Sort use official sort.Sort() for easily implementation
//	Sort sorts data. It makes one call to data.Len to determine n, and O(n*log(n)) calls to data.Less and data.Swap.
//	The sort is not guaranteed to be stable.
func Sort(values []interface{}, comparator Comparator) {
	sort.Sort(sortable{
		values:     values,
		comparator: comparator,
	})
}

// SortStable use official sort.Stable() for easily implementation
//	Stable sorts data while keeping the original order of equal elements.
//	It makes one call to data.Len to determine n, O(n*log(n)) calls to data.Less and O(n*log(n)*log(n)) calls to data.Swap.
func SortStable(values []interface{}, comparator Comparator) {
	sort.Stable(sortable{
		values:     values,
		comparator: comparator,
	})
}

// GetSortedValues returns sorted values
func GetSortedValues(container Container, comparator Comparator, stable bool) []interface{} {
	values := container.Values()
	if len(values) < 2 {
		return values
	}
	if !stable {
		Sort(values, comparator)
	} else {
		SortStable(values, comparator)
	}
	return values
}
