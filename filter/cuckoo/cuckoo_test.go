package cuckoofilter

import (
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"godev/basic"
	"math"
	"testing"
)

func TestNewCuckooFilter(t *testing.T) {
	var _ basic.Container = (*CuckooFilter)(nil)

	cf := NewCuckooFilter(0)
	assert.NotNil(t, cf)
	assert.Equal(t, 1, len(cf.buckets))
	assert.True(t, cf.Empty())

	cf = NewCuckooFilter(128)
	assert.NotNil(t, cf)
	assert.Equal(t, 128/bucketSize, len(cf.buckets))

	cf = NewCuckooFilter(100)
	assert.NotNil(t, cf)
	assert.Equal(t, 128/bucketSize, len(cf.buckets))
}

func TestCuckooFilter_InsertKey(t *testing.T) {
	c := 128
	cf := NewCuckooFilter(uint(c))
	assert.Equal(t, c, len(cf.Values()))
	assert.True(t, cf.Empty())
	var key string
	for i := 0; i < c/2; i++ {
		key = fmt.Sprintf("%s#%d", "foo", i)
		assert.True(t, cf.InsertKey(key))
		assert.True(t, cf.ContainsKey(key))
	}

	assert.False(t, cf.Empty())
	assert.Equal(t, c/2, cf.Size())

	cnt := 0
	for i := 0; i < c/2; i++ {
		key = fmt.Sprintf("%s#%d", "foo", i)
		if cf.ContainsKey(key) {
			cnt++
		}
	}

	assert.Equal(t, c/2, cnt)

	assert.False(t, cf.DeleteKey("foo"))

	for i := 0; i < c/2; i++ {
		key = fmt.Sprintf("%s#%d", "foo", i)
		assert.True(t, cf.DeleteKey(key))
	}

	assert.True(t, cf.Empty())
	assert.Equal(t, 0, cf.Size())

	cf.Clear()
	assert.Equal(t, 0, cf.Size())
}

// this test taken from https://github.com/dgryski/go-cuckoof
// thanks to dgryski
func TestBasicUint32(t *testing.T) {
	loadFactors := []float64{0.25, 0.5, 0.75, 0.95, 0.97, 0.99, 1.25}
	for p := 4; p <= 16; p += 4 {
		for _, lf := range loadFactors {
			size := 1 << uint16(p) // Total capacity
			r := hammer(uint(size), lf)
			// We tried to insert size * lf elements, size * lf * r.fails failed.
			// Thus, effective load is (size * lf - size * lf * r.fails ) / size.
			// size * lf * r.fails elements were kicked out, so the actual
			// false negatives rate is r.failes. Some of those are masked because of false positives.
			effectiveLoad := lf * (1 - r.fails)
			estimatedFalseNegatives := r.fails * (1 - r.falsePositives)
			what := fmt.Sprintf("size: %d(2^%d) load factor: %.00f%% effective load: %0.03f%% %#v efn: %.02f delta: %f", size, p, lf*100, effectiveLoad*100, r, estimatedFalseNegatives, estimatedFalseNegatives-r.falseNegatives)
			if lf < 0.96 {
				// Harold: almost equal, i use murmur3 hash func, dgryski uses metro, this may introduce difference
				if r.fails != 0 || math.Abs(r.falseNegatives) > 0.0001 {
					t.Errorf("Expected failed==0 && falseNegatives==0 --- %s", what)
				}
			}
			if math.Abs(r.falseNegatives-estimatedFalseNegatives) > 0.02 {
				t.Errorf("Expected delta = |falseNegatives - estimatedFalseNegatives| to be small --- %s", what)
			}
			if r.falsePositives > 0.3 {
				t.Errorf("Expected falseNegatives to be less than 0.3 --- %s", what)

			}
			fmt.Println(what)
		}
	}
	return
}

type rates struct {
	fails, falsePositives, falseNegatives float64
}

func hammer(size uint, loadFactor float64) rates {
	f := NewCuckooFilter(size) // bucket size is 4
	num := int(float64(size) * loadFactor)
	elts := make([][]byte, num)
	bts := make([]byte, num*4)
	var r rates
	// Populate the filter.
	for i := range elts {
		b := bts[i*4 : i*4+4]
		binary.BigEndian.PutUint32(b, uint32(i))
		elts[i] = b
		if !f.insert(b) {
			r.fails += 1.0 / float64(len(elts))
		}
	}
	// Check for false negatives.
	for _, b := range elts {
		if !f.lookup(b) {
			r.falseNegatives += 1.0 / float64(len(elts))
		}
	}
	// Check for false positives.
	n := size * 4
	elt := make([]byte, 4)
	for i := 0; i < int(n); i++ {
		binary.BigEndian.PutUint32(elt, uint32(i+num))
		if f.lookup(elt) {
			r.falsePositives += 1.0 / float64(n)
		}
	}
	return r
}
