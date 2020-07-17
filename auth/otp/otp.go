package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"hash"
	"math"
	"strconv"
	"strings"
)

// thanks to https://github.com/pyauth/pyotp
// port from `pyotp`

// HashFunc (sha1/sha256/sha512)
type HashFunc func() hash.Hash

// OTP one time password
type OTP struct {
	Digits int      // number of integers in the OTP. Some apps expect this to be 6 digits, others support more.
	Digest HashFunc // digest function to use in the HMAC (expected to be sha1)
	Secret string
}

// NewOTP ...
func NewOTP(secret string, digits int, hash HashFunc) OTP {
	if hash == nil {
		hash = sha1.New
	}
	if digits < DefaultDigits {
		digits = DefaultDigits
	}
	return OTP{
		Digits: digits,
		Digest: hash,
		Secret: secret,
	}
}

// GenerateOTP generates OTP string
// input: the HMAC counter value to use as the OTP input.
//	Usually either the counter, or the computed integer based on the Unix timestamp (int64)
func (otp *OTP) GenerateOTP(input int64) string {
	if input < 0 {
		panic("input must be positive integer")
	}
	hasher := hmac.New(otp.Digest, otp.byteSecret())
	// message bytes
	bytes, _ := Int64ToBytes(input)
	_, _ = hasher.Write(bytes)
	// MAC
	hmacHash := hasher.Sum(nil)
	code := int(HashToCode(hmacHash)) % int(math.Pow10(otp.Digits))
	return fmt.Sprintf("%0"+strconv.Itoa(otp.Digits)+"d", code)
}

// secret as hash key
func (otp *OTP) byteSecret() []byte {
	// RFC 4648: standard padding char "="
	// StdEncoding: "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	// Any remainder beyond modulus 8 can only be 2, 4, 5, or 7 characters long,
	//	so padding must always be 6, 4, 3 or 1 = character, any other length is invalid.
	missingPadding := len(otp.Secret) % 8
	if missingPadding != 0 {
		otp.Secret += strings.Repeat("=", 8-missingPadding)
	}
	bytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(otp.Secret))
	if err != nil {
		panic(err)
	}
	return bytes
}
