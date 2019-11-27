package arc

import (
	"fmt"
	"goContainer/cache"
	"goContainer/utils"
	"testing"
)

func TestNewARC(t *testing.T) {
	var _ cache.Interface = (*ARC)(nil)
	arc, err := NewARC(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if arc.Size() != 0 {
		t.Fatal(arc.Size(), arc.size)
	}
	k, v, found := arc.GetLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}
	k, v, found = arc.RemoveLeastUsed()
	if k != nil || v != nil || found {
		t.Fatal(k, v, found)
	}

	_, err = NewARC(-10, nil)
	if err == nil || err.Error() != "invalid cache size" {
		t.Fatal(err)
	}
}

func TestARC_Add(t *testing.T) {
	arc, err := NewARC(10, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := arc.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}

	if arc.Size() != 10 {
		t.Fatal(arc.Size(), arc.size)
	}

	for i := range a {
		if !arc.Contains(a[i]) {
			t.Fatal(a[i], i)
		}
	}

	for i := range a {
		found, evicted := arc.Add(a[i], i)
		if !found || evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := arc.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}

	arc.Resize(20)
	for i := range a {
		found, evicted := arc.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}

	arc.Resize(10)
	for i := range b {
		found, evicted := arc.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}

	k, v, found := arc.GetLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	k, v, found = arc.RemoveLeastUsed()
	if k != b[0] || v != 0 || !found {
		t.Fatal(k, b[0], v, 0, found)
	}

	arc.Clear()
	if arc.Size() != 0 {
		t.Fatal(arc.Size(), arc.size)
	}
}

func TestARC_Get(t *testing.T) {
	arc, err := NewARC(20, func(key interface{}, value interface{}) {
		fmt.Printf("key: %v, value: %v pair is deleted\n", key, value)
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	a := make([]string, 10)
	for i := range a {
		a[i] = utils.GenerateRandomString(5)
		found, evicted := arc.Add(a[i], i)
		if found || evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}
	k, v, found := arc.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	b := make([]int, 10)
	for i := range b {
		b[i] = utils.GenerateRandomInt()
		found, evicted := arc.Add(b[i], i)
		if found || evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}
	k, v, found = arc.GetLeastUsed()
	if k != a[0] || v != 0 || !found {
		t.Fatal(k, a[0], v, 0, found)
	}

	for i := range a {
		value, found := arc.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range a {
		value, found := arc.Get(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	fmt.Printf("%v\n", arc.printPairs())

	diff, err := arc.Resize(10)
	if err != nil || diff != -10 {
		t.Fatal(err, diff)
	}

	fmt.Printf("%v\n", arc.printPairs())
	for i := range a {
		value, found := arc.Peek(a[i])
		if !found || value != i {
			t.Fatal(found, a[i], i, value)
		}
	}
	for i := range b {
		value, found := arc.Peek(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}
	for i := range b {
		value, found := arc.Get(b[i])
		if found || value != nil {
			t.Fatal(found, b[i], i, value)
		}
	}

	for i := range b {
		value, found := arc.Remove(b[i])
		if found || value != nil {
			t.Fatal(found, b[i])
		}
	}

	for i := range a {
		value, found := arc.Remove(a[i])
		if !found || value != i {
			t.Fatal(found, a[i])
		}
	}

	for i := range a[:5] {
		found, evicted := arc.Add(a[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}
	for i := range a[:5] {
		value, found := arc.Get(a[i])
		if !found || value != i {
			t.Fatal(found, value, a[i], i)
		}
	}
	for i := range b {
		found, evicted := arc.Add(b[i], i)
		if found || !evicted {
			t.Fatal(found, evicted, arc.Size(), arc.size)
		}
	}
}
