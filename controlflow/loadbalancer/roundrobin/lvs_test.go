package roundrobin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLVS(t *testing.T) {
	wrr := NewLVS()
	err := wrr.SetNodes(nil, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)

	err = wrr.SetNodes(nil, []int{0})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)

	err = wrr.SetNodes([]string{"0.0.0.1"}, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)

	err = wrr.SetNodes([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{2, 4, 8})
	assert.NoError(t, err)

	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
}

func TestLVSrr_Next(t *testing.T) {
	wrr := NewLVS()
	err := wrr.SetNodes([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{0, 0, 0})
	assert.NoError(t, err)

	assert.Empty(t, wrr.Next())
}

func TestLVSrr_Reset(t *testing.T) {
	wrr := NewLVS()
	err := wrr.SetNodes([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{2, 4, 8})
	assert.NoError(t, err)

	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())

	wrr.Reset()
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
}

func TestLVSrr_Distribution(t *testing.T) {
	num := 10000
	nodes := []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}
	weights := []int{2, 4, 8}
	totalWeights := 14
	n := make(map[string]int)
	for i := range nodes {
		n[nodes[i]] = weights[i]
	}
	wrr := NewLVS()
	dist := testWeightedRoundRobinDistribution(t, wrr, nodes, weights, num)
	for k, v := range dist {
		rate := float64(v) / float64(num)
		t.Logf("%s: %0.3f", k, rate)
		if !almostEqual(rate, float64(n[k])/float64(totalWeights), 0.01) {
			t.Fail()
		}
	}
}

func almostEqual(a, b, tolerance float64) bool {
	if a-b < tolerance && b-a < tolerance {
		return true
	}
	return false
}

func testWeightedRoundRobinDistribution(t *testing.T, wrr RoundRobin, nodes []string, weights []int, num int) map[string]int {
	err := wrr.SetNodes(nodes, weights)
	if err != nil || wrr == nil {
		t.Fail()
		return nil
	}
	dist := make(map[string]int)
	for i := 0; i < num; i++ {
		dist[wrr.Next()]++
	}
	return dist
}

func TestGcdArray(t *testing.T) {
	arr := []int{8, 3, 1, 2, 5, 6, 7, 9}
	assert.Equal(t, 1, gcdArray(arr))
}
