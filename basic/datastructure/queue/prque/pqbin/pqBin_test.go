package pqbin

import (
	"fmt"
	"goContainer/basic"
	"goContainer/utils"
	"math"
	"strconv"
	"testing"
)

func TestPQBin(t *testing.T) {
	var _ basic.Container = (*PQBin)(nil)
	pq := NewPQBin()
	if !pq.Empty() {
		t.Fail()
	}

	pq.Push("x", 7)
	if pq.Empty() {
		t.Fail()
	}

	items := map[string]int{
		"a": 6,
		"b": 5,
		"c": 10,
		"d": 3,
		"e": 5,
	}

	for k, v := range items {
		pq.Push(k, v)
	}

	it := pq.Pop().(string)
	if it != "c" {
		fmt.Println(it)
		t.Fail()
	}

	pq.Pop() // pop "x"

	for _, v := range pq.Values() {
		if _, found := items[v.(string)]; !found {
			t.Fail()
		}
	}

	pq.Clear()
	if !pq.Empty() {
		t.Fail()
	}
}

// BenchmarkPQBin_Push-8   	10000000	       253 ns/op
func BenchmarkPQBin_Push(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	b.ResetTimer()

	queue := NewPQBin()
	for i := 0; i < len(data); i++ {
		queue.Push(data[i], data[i])
	}
}

// BenchmarkPQBin_Pop-8   	 1000000	      1493 ns/op
func BenchmarkPQBin_Pop(b *testing.B) {
	data := make([]int, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = utils.GenerateRandomInt()
	}
	queue := NewPQBin()
	for i := 0; i < len(data); i++ {
		queue.Push(data[i], data[i])
	}
	b.ResetTimer()

	for i := 0; i < len(data); i++ {
		queue.Pop()
	}
}

func BenchmarkPQBin(b *testing.B) {
	for k := 1.0; k <= 5; k++ {
		n := int(math.Pow(10, k))

		pq := NewPQBin()
		rn := 0
		for i := 0; i < n; i++ {
			rn = utils.GenerateRandomInt()
			pq.Push(rn, rn)
		}
		num := utils.GenerateRandomInt()
		b.ResetTimer()

		b.Run("PQBin Push: size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				pq.Pop()
				pq.Push(num, num)
			}
		})
	}
}

/*
BenchmarkPQBin/PQBin_Push:_size-10-8         	20000000	       105 ns/op
BenchmarkPQBin/PQBin_Push:_size-100-8        	20000000	       107 ns/op
BenchmarkPQBin/PQBin_Push:_size-1000-8       	10000000	       112 ns/op
BenchmarkPQBin/PQBin_Push:_size-10000-8      	20000000	       116 ns/op
BenchmarkPQBin/PQBin_Push:_size-100000-8     	10000000	       129 ns/op
*/
