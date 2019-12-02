package tmpobjpool

import "testing"

func TestPool(t *testing.T) {
	p := New()

	p.AddSize(10)
	
	if x := p.Get(10); x != nil {
		t.Fatal(x)
	}

	p.Put("asd", 10)
	if x := p.Get(10); x != "asd" {
		t.Fatal(x)
	}

	p.Put(123, 5)
	if x := p.Get(5); x.(int) != 123 {
		t.Fatal(x)
	}
}