package loadbalancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRRW(t *testing.T) {
	rrw, err := NewRRW(nil, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)
	assert.Nil(t, rrw)

	rrw, err = NewRRW(nil, []int{0})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)
	assert.Nil(t, rrw)

	rrw, err = NewRRW([]string{"0.0.0.1"}, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)
	assert.Nil(t, rrw)

	rrw, err = NewRRW([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{2, 4, 8})
	assert.NoError(t, err)
	assert.NotNil(t, rrw)

	assert.Equal(t, "0.0.0.3", rrw.Next())
}

func TestRoundRobinWeighted_Next(t *testing.T) {
	rrw, err := NewRRW([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{0, 0, 0})
	assert.NoError(t, err)
	assert.NotNil(t, rrw)

	assert.Empty(t, rrw.Next())
}

func TestRoundRobinWeighted_Distribution(t *testing.T) {
	num := 10000
	nodes := []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}
	weights := []int{2, 4, 8}
	totalWeights := 14
	n := make(map[string]int)
	for i := range nodes {
		n[nodes[i]] = weights[i]
	}
	dist := testRoundRobinWeightedDistribution(t, nodes, weights, num)
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

func testRoundRobinWeightedDistribution(t *testing.T, nodes []string, weights []int, num int) map[string]int {
	rrw, err := NewRRW(nodes, weights)
	if err != nil || rrw == nil {
		t.Fail()
		return nil
	}
	dist := make(map[string]int)
	for i := 0; i < num; i++ {
		dist[rrw.Next()]++
	}
	return dist
}

func TestGcdArray(t *testing.T) {
	arr := []int{8, 3, 1, 2, 5, 6, 7, 9}
	assert.Equal(t, 1, gcdArray(arr))
}
