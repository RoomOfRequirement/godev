package roundrobin

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New(Simple)
	assert.NoError(t, err)
	_, err = New(LVS)
	assert.NoError(t, err)
	_, err = New(Nginx)
	assert.NoError(t, err)

	_, err = New(100)
	assert.Error(t, err)
	assert.Equal(t, ErrUnsupportedAlgorithm, err)
}

func BenchmarkNext(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	nodes := make([]string, 0, 10)
	weights := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, strconv.Itoa(i))
		weights = append(weights, rand.Intn(100))
	}

	b.Run("Simple", func(b *testing.B) {
		b.ReportAllocs()
		rr, _ := New(Simple)
		_ = rr.SetNodes(nodes, nil)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = rr.Next()
		}
	})

	b.Run("LVS", func(b *testing.B) {
		b.ReportAllocs()
		rr, _ := New(LVS)
		_ = rr.SetNodes(nodes, weights)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = rr.Next()
		}
	})

	b.Run("Nginx", func(b *testing.B) {
		b.ReportAllocs()
		rr, _ := New(Nginx)
		_ = rr.SetNodes(nodes, weights)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = rr.Next()
		}
	})
}
