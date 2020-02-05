package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase64(t *testing.T) {
	src := "hello world"
	str64 := Base64Encode(src)
	dbytes, err := Base64Decode(str64)
	assert.NoError(t, err)
	assert.Equal(t, src, BytesToString(dbytes))
}
