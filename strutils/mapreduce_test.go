package strutils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	s := "hello world nice day feel good play hard"
	out := Map(strings.Split(s, " "), func(s string) string {
		return s + "|" + s
	})
	assert.Equal(t, strings.Split("hello|hello world|world nice|nice day|day feel|feel good|good play|play hard|hard", " "), out)
}

func TestReduce(t *testing.T) {
	s := "hello world nice day feel good play hard"
	out := Reduce(Map(strings.Split(s, " "), func(s string) string {
		return s + "|" + s
	}), func(strs []string) string {
		return strings.Join(strs, " ")
	})
	assert.Equal(t, "hello|hello world|world nice|nice day|day feel|feel good|good play|play hard|hard", out)
}

func TestMapReduce(t *testing.T) {
	s := "hello world nice day feel good play hard"
	out := MapReduce(strings.Split(s, " "), func(s string) string {
		return s + "|" + s
	}, func(strs []string) string {
		return strings.Join(strs, " ")
	})
	assert.Equal(t, "hello|hello world|world nice|nice day|day feel|feel good|good play|play hard|hard", out)
}

func TestSplitMapReduce(t *testing.T) {
	s := "hello world nice day feel good play hard"
	out := SplitMapReduce(s, " ", func(s string) string {
		return s + "|" + s
	}, func(strs []string) string {
		return strings.Join(strs, " ")
	})

	assert.Equal(t, "hello|hello world|world nice|nice day|day feel|feel good|good play|play hard|hard", out)
}
