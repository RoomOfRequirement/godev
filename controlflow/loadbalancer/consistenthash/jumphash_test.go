package consistenthash

import (
	"encoding/base64"
	"math/rand"
	"testing"
)

func TestJumpHash(t *testing.T) {
	virtual := DefaultVirtual
	nodes := map[int]string{
		0: "0.0.0.0", // [0 - virtual)
		1: "0.0.0.1", // [virtual - 2 * virtual)
		2: "0.0.0.2", // [2 * virtual - 3 * virtual)
	}
	numBuckets := len(nodes) * virtual
	var hashF HashFunc = hash
	dist := map[string]int{
		"0.0.0.0": 0,
		"0.0.0.1": 0,
		"0.0.0.2": 0,
	}

	numKey := 100000
	buf := make([]byte, 12)
	for i := 0; i < numKey; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		h := hashF([]byte(base64.StdEncoding.EncodeToString(buf)))
		n := JumpHash(uint64(h), numBuckets)
		switch n / int32(virtual) {
		case 0:
			dist[nodes[0]]++
		case 1:
			dist[nodes[1]]++
		case 2:
			dist[nodes[2]]++
		}
	}
	ratios := make([]float64, 0, len(nodes))
	for k, v := range dist {
		ratio := float64(v) / float64(numKey)
		ratios = append(ratios, ratio)
		t.Logf("%s: %0.3f", "jump hash: "+k, ratio)
	}
	for i := 0; i < len(ratios)-1; i++ {
		if !almostEqual(ratios[i], ratios[i+1], 0.01) {
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
