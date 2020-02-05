package strutils

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestShuffle(t *testing.T) {
	s := "test a, b, c d"
	ss := Shuffle(s, rand.NewSource(time.Now().Unix()))
	assert.NotEqual(t, s, ss)
	t.Log(ss)
}
