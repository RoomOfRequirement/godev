package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructToString(t *testing.T) {
	s := struct {
		S string
		B bool
		N int
		F float32
		M map[string]string
	}{
		"s",
		true,
		1,
		0.2,
		map[string]string{
			"a": "a",
		},
	}
	str, err := StructToString(s)
	assert.NoError(t, err)
	assert.Equal(t, "{\n\t\"S\": \"s\",\n\t\"B\": true,\n\t\"N\": 1,\n\t\"F\": 0.2,\n\t\"M\": {\n\t\t\"a\": \"a\"\n\t}\n}", str)

	str, err = StructToString(func() {})
	assert.Error(t, err)
	assert.Equal(t, "", str)
}

func TestStringToStruct(t *testing.T) {
	str := "{\"S\": \"s\",\n\t\"B\": true,\n\t\"N\": 1,\n\t\"F\": 0.2,\n\t\"M\": {\n\t\t\"a\": \"a\"\n\t}}"
	s := &struct {
		S string
		B bool
		N int
		F float32
		M map[string]string
	}{}
	err := StringToStruct(str, s)
	assert.NoError(t, err)
	assert.EqualValues(t, struct {
		S string
		B bool
		N int
		F float32
		M map[string]string
	}{
		"s",
		true,
		1,
		0.2,
		map[string]string{
			"a": "a",
		},
	}, *s)

	err = StringToStruct(str, nil)
	assert.Error(t, err)
}
