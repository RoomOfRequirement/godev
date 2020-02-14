package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"strconv"
	"testing"
)

func BenchmarkGenerateRandomInt(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				for j := 0; j < n; j++ {
					_ = GenerateRandomInt()
				}
			}
		})
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				_ = GenerateRandomString(n)
			}
		})
	}
}

func TestGenerateRandomInt(t *testing.T) {
	var n interface{}
	n = GenerateRandomInt()
	_, ok := n.(int)
	if !ok {
		t.Fail()
	}
}

func TestGenerateRandomIntInRange(t *testing.T) {
	s := GenerateRandomInt()
	e := s + 1000
	n := GenerateRandomIntInRange(s, e)
	if n < s || n >= e {
		t.Fail()
	}

	n = GenerateRandomIntInRange(e, s)
	if n < s || n >= e {
		t.Fail()
	}
}

func TestGenerateRandomString(t *testing.T) {
	var s interface{}
	s = GenerateRandomString(10)
	str, ok := s.(string)
	if !ok || len(str) != 10 {
		t.Fail()
	}
}

func TestPartialFunc(t *testing.T) {
	funcMap := map[string]interface{}{
		"f1": func(a int, b string) string {
			return strconv.Itoa(a) + b
		},
		"f2": func(m map[string]string) string {
			res := ""
			for k, v := range m {
				res += k + v
			}
			return res
		},
	}

	resSlice, err := PartialFunc(funcMap, "f1", 10)
	assert.Error(t, err)

	resSlice, err = PartialFunc(funcMap, "f1", 10, "a")
	assert.NoError(t, err)
	assert.Equal(t, "10a", resSlice[0].String())
	resSlice, err = PartialFunc(funcMap, "f2", map[string]string{
		"a": "1",
		"b": "2",
	})
	assert.NoError(t, err)
	assert.True(t, resSlice[0].String() == "a1b2" || resSlice[0].String() == "b2a1")
}

func TestTrace(t *testing.T) {
	panicFunc := func() interface{} {
		panic("hi, panic here, need help")
		return nil
	}
	call := func() interface{} {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Println(Trace(message))
			}
		}()
		return panicFunc()
	}
	assert.NotPanics(t, func() {
		call()
	})
}
