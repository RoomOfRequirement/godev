package strutils

import (
	"unicode"
	"unicode/utf8"
)

// IsAllASCII returns true if all runes of input string is ascii
func IsAllASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// from strings.Fields
// https://golang.org/src/strings/strings.go?s=8396:8426#L319
func isAllASCII(s string) bool {
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(s); i++ {
		setBits |= s[i]
	}

	return setBits < utf8.RuneSelf
}
