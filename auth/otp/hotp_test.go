package otp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strconv"
	"testing"
)

func TestHOTP_At(t *testing.T) {
	hotp := NewHOTP("base32secret3232", 6, "SHA1", nil)
	assert.Equal(t, "260182", hotp.At(0))
	assert.Equal(t, "055283", hotp.At(1))
	assert.Equal(t, "316439", hotp.At(1401))
}

func TestHOTP_Verify(t *testing.T) {
	hotp := NewHOTP("base32secret3232", 6, "SHA1", nil)
	assert.True(t, hotp.Verify("260182", 0))
	assert.False(t, hotp.Verify("260182", 1))
}

func TestHOTP_ProvisioningUri(t *testing.T) {
	issuer := "example"
	name := "test"
	secret, _ := GenerateRandomSecret(16)
	counter := 10
	expected, _ := url.Parse(fmt.Sprintf("otpauth://%s/%s", Hotp, issuer+":"+name))
	q := expected.Query()
	q.Set("secret", secret)
	q.Set("counter", strconv.Itoa(counter))
	q.Set("issuer", issuer)
	expected.RawQuery = q.Encode()
	hotp := NewHOTP(secret, 0, "", nil)
	got, err := url.Parse(hotp.ProvisioningURI(name, issuer, counter))
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
