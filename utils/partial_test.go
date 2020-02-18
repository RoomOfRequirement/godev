package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

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
