package roundrobin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNgx(t *testing.T) {
	wrr := NewNgx()
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
	assert.Equal(t, "0.0.0.2", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
}

func TestNgx_Next(t *testing.T) {
	wrr := NewNgx()
	err := wrr.SetNodes([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{0, 0, 0})
	assert.NoError(t, err)

	assert.Empty(t, wrr.Next())
}

func TestNgx_Reset(t *testing.T) {
	wrr := NewNgx()
	err := wrr.SetNodes([]string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, []int{5, 1, 1})
	assert.NoError(t, err)

	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())

	wrr.Reset()
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.2", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.3", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
	assert.Equal(t, "0.0.0.1", wrr.Next())
}

func TestNgx_Distribution(t *testing.T) {
	num := 10000
	nodes := []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}
	weights := []int{2, 4, 8}
	totalWeights := 14
	n := make(map[string]int)
	for i := range nodes {
		n[nodes[i]] = weights[i]
	}
	wrr := NewNgx()
	dist := testWeightedRoundRobinDistribution(t, wrr, nodes, weights, num)
	for k, v := range dist {
		rate := float64(v) / float64(num)
		t.Logf("%s: %0.3f", k, rate)
		if !almostEqual(rate, float64(n[k])/float64(totalWeights), 0.01) {
			t.Fail()
		}
	}
}
