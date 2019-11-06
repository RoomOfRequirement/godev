package bag

import (
	"bytes"
	"fmt"
)

// Bag struct / Multiset
//	https://en.wikipedia.org/wiki/Set_(abstract_data_type)#Multiset
type Bag struct {
	m    map[interface{}]int
	size int
}

// Entry struct
type Entry struct {
	entry interface{}
	count int
}

// GetEntry returns inside entry
func (e *Entry) GetEntry() interface{} {
	return e.entry
}

// GetCount returns inside entry's quantity
func (e *Entry) GetCount() int {
	return e.count
}

// NewBag creates a new bag
func NewBag() *Bag {
	return &Bag{
		m:    make(map[interface{}]int),
		size: 0,
	}
}

// Add adds entry into bag
func (bag *Bag) Add(x interface{}) {
	if _, found := bag.m[x]; found {
		bag.m[x]++
	} else {
		bag.m[x] = 1
	}
	bag.size++
}

// Count returns entry's quantity
func (bag *Bag) Count(x interface{}) int {
	if _, found := bag.m[x]; found {
		return bag.m[x]
	}
	return 0
}

// Contains returns true if entry inside the bag
func (bag *Bag) Contains(x interface{}) bool {
	return bag.Count(x) != 0
}

// SetCount sets entry with quantity
func (bag *Bag) SetCount(x interface{}, cnt int) {
	bag.m[x] = cnt
}

// DeleteAll deletes all same entries
func (bag *Bag) DeleteAll(x interface{}) bool {
	if cnt, found := bag.m[x]; found {
		delete(bag.m, x)
		bag.size -= cnt
		return true
	}
	return false
}

// DeleteOne deletes entry quantity by one
func (bag *Bag) DeleteOne(x interface{}) bool {
	if cnt, found := bag.m[x]; found {
		if cnt > 1 {
			bag.m[x]--
		} else {
			delete(bag.m, x)
		}
		bag.size--
		return true
	}
	return false
}

// Entries returns all distinct entries inside the bag
func (bag *Bag) Entries() []interface{} {
	es := make([]interface{}, 0, len(bag.m))
	for k := range bag.m {
		es = append(es, k)
	}
	return es
}

// EntriesWithCount returns entries with count as Entry type
func (bag *Bag) EntriesWithCount() []Entry {
	entries := make([]Entry, 0, len(bag.m))
	for k, v := range bag.m {
		entries = append(entries, Entry{
			entry: k,
			count: v,
		})
	}
	return entries
}

// ForEachEntry maps function on each entry
//	notice: for entry has n quantity, func f will executes n times on this entry
func (bag *Bag) ForEachEntry(f func(interface{})) {
	for k, count := range bag.m {
		for i := 0; i < count; i++ {
			f(k)
		}
	}
}

// Merge merges two bags
//	notice: inplace change
func (bag *Bag) Merge(other *Bag) {
	for k, v := range other.m {
		if bag.Contains(k) {
			bag.m[k] += v
		} else {
			bag.m[k] = v
		}

		bag.size += v
	}
}

// ContainsAll returns true if current bag contains all entries of other bag (with at least quantity one)
func (bag *Bag) ContainsAll(other *Bag) bool {
	for k := range other.m {
		if !bag.Contains(k) {
			return false
		}
	}
	return true
}

// String for pretty print
func (bag *Bag) String() string {
	var buf bytes.Buffer
	buf.WriteString("{ ")

	bag.ForEachEntry(func(x interface{}) {
		_, _ = fmt.Fprintf(&buf, "%v ", x)
	})

	buf.WriteString("}\n")
	return buf.String()
}

// Empty returns true if no entries inside the bag
func (bag *Bag) Empty() bool {
	return bag.size == 0
}

// Size returns quantity of entries inside the bag
func (bag *Bag) Size() int {
	return bag.size
}

// Clear clears bag
func (bag *Bag) Clear() {
	*bag = *NewBag()
}

// Values returns all entries inside the bag
func (bag *Bag) Values() []interface{} {
	return bag.Entries()
}
