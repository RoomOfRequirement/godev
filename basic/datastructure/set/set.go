package set

type void struct{}

// IntSet int set
type IntSet map[int]void

// NewIntSet returns a new IntSet with capacity of `capacity`
func NewIntSet(capacity int) *IntSet {
	if capacity < 0 {
		capacity = 0
	}
	s := make(IntSet, capacity)
	return &s
}

// NewIntSetFromIntArray returns a new IntSet from int array
func NewIntSetFromIntArray(array []int) *IntSet {
	s := make(IntSet, len(array))
	for _, i := range array {
		s[i] = void{}
	}
	return &s
}

// IntSetEqual returns true if two IntSet have same values inside and have same size
func IntSetEqual(s, s1 *IntSet) bool {
	if s.Size() != s1.Size() || s.Difference(s1).Size() != 0 {
		return false
	}
	return true
}

// Empty return true if IntSet size equals to 0
func (s *IntSet) Empty() bool {
	return s.Size() == 0
}

// Clear clears IntSet by assigning a new IntSet with the same size to its pointer
func (s *IntSet) Clear() {
	*s = *(NewIntSet(s.Size()))
}

// Values returns values stored inside the set
func (s *IntSet) Values() []interface{} {
	keys := make([]interface{}, 0, s.Size())
	for k := range *s {
		keys = append(keys, k)
	}
	return keys
}

// Get returns true if input `x` found in the set
func (s *IntSet) Get(x int) bool {
	_, ok := (*s)[x]
	return ok
}

// Add adds `x` into the set
func (s *IntSet) Add(x int) {
	(*s)[x] = void{}
}

// Delete deletes `x` from the set
func (s *IntSet) Delete(x int) {
	delete(*s, x)
}

// Size returns the number of ints inside the set
func (s *IntSet) Size() int {
	return len(*s)
}

// Union returns a new set which combines two sets
func (s *IntSet) Union(s1 *IntSet) *IntSet {
	ns := make(IntSet, s.Size()+s1.Size())
	for k := range *s {
		ns[k] = void{}
	}
	for k := range *s1 {
		ns[k] = void{}
	}
	return &ns
}

// Intersection returns a new set with common values from two sets
func (s *IntSet) Intersection(s1 *IntSet) *IntSet {
	tmp0, tmp1 := s, s1
	if s.Size() > s1.Size() {
		tmp0 = s1
		tmp1 = s
	}
	ns := make(IntSet, tmp0.Size())
	for k := range *tmp0 {
		if _, ok := (*tmp1)[k]; ok {
			ns[k] = void{}
		}
	}
	return &ns
}

// Difference returns a new set with different values from two sets
func (s *IntSet) Difference(s1 *IntSet) *IntSet {
	ns := s.Union(s1)
	for k := range *(s.Intersection(s1)) {
		ns.Delete(k)
	}
	return ns
}

// FloatSet float set
type FloatSet map[float64]void

// NewFloatSet returns a new set consists of float64s
func NewFloatSet(capacity int) *FloatSet {
	if capacity < 0 {
		capacity = 0
	}
	s := make(FloatSet, capacity)
	return &s
}

// NewFloatSetFromFloatArray returns a new set from float64 array
func NewFloatSetFromFloatArray(array []float64) *FloatSet {
	s := make(FloatSet, len(array))
	for _, i := range array {
		s[i] = void{}
	}
	return &s
}

// FloatSetEqual returns true if two float sets have same values and same size
func FloatSetEqual(s, s1 *FloatSet) bool {
	if s.Size() != s1.Size() || s.Difference(s1).Size() != 0 {
		return false
	}
	return true
}

// Empty returns true if size of FloatSet equals to 0
func (s *FloatSet) Empty() bool {
	return s.Size() == 0
}

// Clear clears IntSet by assigning a new FloatSet with the same size to its pointer
func (s *FloatSet) Clear() {
	*s = *(NewFloatSet(s.Size()))
}

// Values returns values stored inside the set
func (s *FloatSet) Values() []interface{} {
	keys := make([]interface{}, 0, s.Size())
	for k := range *s {
		keys = append(keys, k)
	}
	return keys
}

// Get returns true if input `x` found in the set
func (s *FloatSet) Get(x float64) bool {
	_, ok := (*s)[x]
	return ok
}

// Add adds `x` into the set
func (s *FloatSet) Add(x float64) {
	(*s)[x] = void{}
}

// Delete deletes `x` from the set
func (s *FloatSet) Delete(x float64) {
	delete(*s, x)
}

// Size returns the number of ints inside the set
func (s *FloatSet) Size() int {
	return len(*s)
}

// Union returns a new set which combines two sets
func (s *FloatSet) Union(s1 *FloatSet) *FloatSet {
	ns := make(FloatSet, s.Size()+s1.Size())
	for k := range *s {
		ns[k] = void{}
	}
	for k := range *s1 {
		ns[k] = void{}
	}
	return &ns
}

// Intersection returns a new set with common values from two sets
func (s *FloatSet) Intersection(s1 *FloatSet) *FloatSet {
	tmp0, tmp1 := s, s1
	if s.Size() > s1.Size() {
		tmp0 = s1
		tmp1 = s
	}
	ns := make(FloatSet, tmp0.Size())
	for k := range *tmp0 {
		if _, ok := (*tmp1)[k]; ok {
			ns[k] = void{}
		}
	}
	return &ns
}

// Difference returns a new set with different values from two sets
func (s *FloatSet) Difference(s1 *FloatSet) *FloatSet {
	ns := s.Union(s1)
	for k := range *(s.Intersection(s1)) {
		ns.Delete(k)
	}
	return ns
}
