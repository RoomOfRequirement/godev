package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	str := "hello world"
	data, err := Compress(str)
	assert.NoError(t, err)

	ret, err := Decompress(data)
	assert.NoError(t, err)

	assert.Equal(t, str, ret)
}
