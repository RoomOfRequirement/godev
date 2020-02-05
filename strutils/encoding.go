package strutils

import "encoding/base64"

// Base64Encode encodes src bytes into base64 string
func Base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString(StringToBytes(src))
}

// Base64Decode decodes str into dst bytes
func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
