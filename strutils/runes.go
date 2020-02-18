package strutils

import (
	"unicode"
	"unicode/utf8"
)

// IsMark returns true when input rune is a mark
//	this is extracted from reverse.go
func IsMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// RunesToBytes converts rune slices to byte slice
func RunesToBytes(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

// IsLetter ...
func IsLetter(r rune) bool {
	return unicode.IsLetter(r)
}

// IsNumber ...
func IsNumber(r rune) bool {
	return unicode.IsNumber(r)
}

// IsLetterOrNumber ...
func IsLetterOrNumber(r rune) bool {
	return IsLetter(r) || IsNumber(r)
}

// IsASCII ...
func IsASCII(r rune) bool {
	return r < unicode.MaxASCII
}
