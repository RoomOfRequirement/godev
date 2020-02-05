package strutils

import "unicode"

// IsMark returns true when input rune is a mark
//	this is extracted from reverse.go
func IsMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}
