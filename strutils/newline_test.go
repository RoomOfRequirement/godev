package strutils

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestNewLine(t *testing.T) {
	if runtime.GOOS == "windows" {
		assert.Equal(t, "\r\n", NewLine())
	}
	assert.Equal(t, "\n", NewLine())
}
