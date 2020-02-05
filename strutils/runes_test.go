package strutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMark(t *testing.T) {
	assert.False(t, IsMark([]rune("there exists no mark rune")[0]))
	assert.False(t, IsMark([]rune(" ")[0]))
	assert.True(t, IsMark([]rune("ò")[1]))
}
