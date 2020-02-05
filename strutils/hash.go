package strutils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"hash/fnv"
)

// HashString returns hash string of input string by input algorithm
func HashString(str, algo string) (string, error) {
	var method hash.Hash
	switch algo {
	case "md5":
		method = md5.New()
	case "sha1":
		method = sha1.New()
	case "sha256":
		method = sha256.New()
	case "sha512":
		method = sha512.New()
	case "fnv32":
		method = fnv.New32()
	case "fnv32a":
		method = fnv.New32a()
	case "fnv64":
		method = fnv.New64()
	case "fnv64a":
		method = fnv.New64a()
	case "fnv128":
		method = fnv.New128()
	case "fnv128a":
		method = fnv.New128a()
	default:
		return "", errors.New("unsupported hash algorithm")
	}
	method.Write(StringToBytes(str))
	return hex.EncodeToString(method.Sum(nil)), nil
}
