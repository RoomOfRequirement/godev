package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashString(t *testing.T) {
	src := "hello world"
	algo := []string{"md5", "sha1", "sha256", "sha512", "fnv32", "fnv32a", "fnv64", "fnv64a", "fnv128", "fnv128a"}
	for _, a := range algo {
		_, err := HashString(src, a)
		assert.NoError(t, err)
	}

	_, err := HashString(src, "test")
	assert.Error(t, err)
}
