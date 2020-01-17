package bloomfilter

import (
	"fmt"
	"godev/basic"
	"godev/utils"
	"strconv"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	var _ basic.Container = (*BloomFilter)(nil)
	n, p := 1000, 0.01
	bf := NewBloomFilter(uint(n), p)
	if bf.Empty() {
		t.Fail()
	}
	for i := 0; i < n; i++ {
		bf.InsertKey(strconv.Itoa(i))
	}

	cnt := 0
	for i := 0; i < n; i++ {
		if !bf.ContainsKey(strconv.Itoa(i)) || bf.ContainsKey(utils.GenerateRandomString(6)) {
			cnt++
		}
	}
	errRate := float64(cnt) / float64(n)
	if errRate > p*2 {
		t.Fatal(fmt.Sprintf("size: %d, supposed error rate: %f, real error rate: %f", n, p, errRate))
	}

	if bf.Size() != len(bf.Values()) {
		t.Fail()
	}

	bf.Clear()
	if !bf.Empty() || bf.Size() != 0 {
		t.Fail()
	}
}
