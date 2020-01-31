package strutils

import (
	"unicode"
	"unicode/utf8"
)

// Reverse ...
//	notice: not work on combining characters (e.g. "bròwn")
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ReversePreservingCombiningCharacters ...
//	notice: code from here: http://rosettacode.org/wiki/Reverse_a_string#Go
//	some discussion found here: https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go/20225618#20225618
func ReversePreservingCombiningCharacters(s string) string {
	if s == "" {
		return ""
	}
	p := []rune(s)
	r := make([]rune, len(p))
	start := len(r)
	for i := 0; i < len(p); {
		// quietly skip invalid UTF-8
		if p[i] == utf8.RuneError {
			i++
			continue
		}
		j := i + 1
		for j < len(p) && (unicode.Is(unicode.Mn, p[j]) ||
			unicode.Is(unicode.Me, p[j]) || unicode.Is(unicode.Mc, p[j])) {
			j++
		}
		for k := j - 1; k >= i; k-- {
			start--
			r[start] = p[k]
		}
		i = j
	}
	return string(r[start:])
}
