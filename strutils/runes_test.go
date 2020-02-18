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

func TestRunesToBytes(t *testing.T) {
	bytes := RunesToBytes([]rune("ò"))
	assert.Equal(t, []byte("ò"), bytes)
}

func TestIsLetter(t *testing.T) {
	assert.True(t, IsLetter([]rune("a")[0]))
	assert.False(t, IsLetter([]rune("ò")[1]))
	assert.False(t, IsLetter([]rune("0")[0]))
}

func TestIsNumber(t *testing.T) {
	assert.True(t, IsNumber([]rune("0")[0]))
	assert.False(t, IsNumber([]rune("a")[0]))
	assert.False(t, IsNumber([]rune("ò")[1]))
}

func TestIsLetterOrNumber(t *testing.T) {
	assert.True(t, IsLetterOrNumber([]rune("a")[0]))
	assert.True(t, IsLetterOrNumber([]rune("0")[0]))
	assert.False(t, IsLetterOrNumber([]rune("ò")[1]))
}

func TestIsASCII(t *testing.T) {
	assert.True(t, IsASCII([]rune("a")[0]))
	assert.True(t, IsASCII([]rune("0")[0]))
	assert.True(t, IsASCII([]rune("&")[0]))
	assert.False(t, IsASCII([]rune("ò")[1]))
}
