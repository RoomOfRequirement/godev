package heap

// Heap interface
type Heap interface {
	FindMin() Item
	DeleteMin() Item
	Insert(item Item)

	// Delete delete item from heap and return it
	// Delete(item Item) Item

	// Meld combine two heaps into one
	// Meld(h Heap) Heap

	// Following to meet `container.Container`
	Size() int
	Empty() bool
	Clear()
	Values() []interface{}
}

// Item interface stands for item stored inside heap
type Item interface {
	// like container.Comparator
	//	 1 if a > b
	//	 0 if a == b
	//	-1 if a < b
	Compare(item Item) int
}
