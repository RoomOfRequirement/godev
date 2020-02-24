package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTernary(t *testing.T) {
	a, b := -1, 1
	assert.Equal(t, a, Ternary(a < b, a, b).(int))
	assert.Equal(t, b, Ternary(a > b, a, b).(int))
}
