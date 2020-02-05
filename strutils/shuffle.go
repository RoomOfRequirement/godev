package strutils

import (
	"math/rand"
	"strings"
)

// Shuffle shuffles input string
func Shuffle(str string, rSrc rand.Source) string {
	r := rand.New(rSrc)
	words := strings.Fields(str)
	r.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return strings.Join(words, " ")
}
