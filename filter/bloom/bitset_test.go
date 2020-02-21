package bloomfilter

import "testing"

func TestBitSet(t *testing.T) {
	bs := NewBitSet(1)
	bs.SetOne(12)
	if !bs.IsOne(12) {
		t.Fail()
	}
	bs.SetZero(12)
	if bs.IsOne(12) {
		t.Fail()
	}
}
