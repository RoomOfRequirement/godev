package otp

import (
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOTP_GenerateOTP(t *testing.T) {
	otp := NewOTP("abc", 8, sha256.New)
	assert.Panics(t, func() {
		otp.GenerateOTP(-1)
	})
	assert.Panics(t, func() {
		otp.byteSecret()
	})
}
