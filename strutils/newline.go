package strutils

import "runtime"

// NewLine ...
func NewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
