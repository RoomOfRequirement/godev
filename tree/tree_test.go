package tree

import (
	"goContainer"
	"goContainer/heap/bheap"
	"goContainer/tree/btree"
	"goContainer/tree/rbtree"
	"goContainer/utils"
	"testing"
)

/*
// this benchmark way is wrong!!!
func BenchmarkTree_InsertOne(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		rbTree := new(rbtree.RBTree)
		rbTree.Comparator = container.IntComparator

		minH := new(bheap.MinHeap)
		minH.Comparator = container.IntComparator
		minH.Init()

		bTree := btree.NewBTree(10, container.IntComparator)

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()

			rbTree.Insert(rn)

			minH.Push(rn)

			bTree.Insert(&btree.Item{
				Key:   rn,
				Value: rn,
			})
		}

		num := container.GenerateRandomInt()
		b.ResetTimer()

		b.Run("Red-Black Tree: size-"+strconv.Itoa(n), func(b *testing.B) {

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				rbTree.Insert(num)
			}
		})

		b.Run("Binary-Heap: size-"+strconv.Itoa(n), func(b *testing.B) {

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				minH.Push(num)
			}
		})

		b.Run("B Tree: size-"+strconv.Itoa(n), func(b *testing.B) {

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				bTree.Insert(&btree.Item{
					Key:   num,
					Value: num,
				})
			}
		})
		fmt.Println(rbTree.Size(), minH.Size(), bTree.Size())
		// only bTree's size is as expected, this is due to its `insertInternal` method,
		// which will not increase item num with the same key
	}
}
*/

// the following should be correct but not as expected....

func BenchmarkTree_InsertOne(b *testing.B) {
	b.Run("Red-Black Tree: ", func(b *testing.B) {
		rbTree := new(rbtree.RBTree)
		rbTree.Comparator = container.IntComparator
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbTree.Insert(data[i])
		}
	})

	b.Run("Binary-Heap: ", func(b *testing.B) {
		minH := new(bheap.MinHeap)
		minH.Comparator = container.IntComparator
		minH.Init()
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			minH.Push(data[i])
		}
	})

	b.Run("B Tree: ", func(b *testing.B) {
		bTree := btree.NewBTree(10, container.IntComparator)
		data := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			data[i] = utils.GenerateRandomInt()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bTree.Insert(&btree.Item{
				Key:   data[i],
				Value: data[i],
			})
		}
	})
}

// BenchmarkTree_InsertOne/Red-Black_Tree:_-8         	 1000000	      1475 ns/op
// BenchmarkTree_InsertOne/Binary-Heap:_-8            	10000000	       139 ns/op
// BenchmarkTree_InsertOne/B_Tree:_-8                 	 1000000	      1560 ns/op
