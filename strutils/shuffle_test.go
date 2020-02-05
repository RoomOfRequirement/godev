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
	// t.Log(ss)
}

func TestShuffleStrs(t *testing.T) {
	strs := []string {"hello world", "good day", "how are you", "fine, thanks", "what's up", "nothing new"}
	sstrs := ShuffleStrs(strs, rand.NewSource(time.Now().Unix()))
	assert.NotEqual(t, strs, sstrs)
	// t.Log(strs)
	// t.Log(sstrs)
}
