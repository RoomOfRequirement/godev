package otp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestTOTP_At(t *testing.T) {
	totp := NewTOTP("base32secret3232", 0, 0, "", nil)
	assert.Equal(t, totp.Now(), totp.At(time.Now().Unix(), 0))
}

func TestTOTP_NowWithExpiration(t *testing.T) {
	totp := NewTOTP("base32secret3232", 0, 0, "", nil)
	otp, exp := totp.NowWithExpiration()
	now := time.Now()
	assert.Equal(t, totp.At(now.Unix(), 0), otp)
	assert.Equal(t, now.Add(time.Duration(totp.Interval)*time.Second).Unix(), exp.Unix())
}

func TestTOTP_Verify(t *testing.T) {
	totp := NewTOTP("base32secret3232", 0, 0, "", nil)
	assert.True(t, totp.Verify(totp.Now(), time.Now().Unix(), 0))
	assert.False(t, totp.Verify(totp.Now(), time.Now().Add(time.Duration(totp.Interval*2)*time.Second).Unix(), 1))
	assert.True(t, totp.Verify(totp.Now(), time.Now().Add(time.Duration(totp.Interval*2)*time.Second).Unix(), 3))
}

func TestTOTP_ProvisioningUri(t *testing.T) {
	issuer := "example"
	name := "test"
	secret, _ := GenerateRandomSecret(16)
	period := 60
	expected, _ := url.Parse(fmt.Sprintf("otpauth://%s/%s", Totp, issuer+":"+name))
	q := expected.Query()
	q.Set("secret", secret)
	q.Set("period", strconv.Itoa(period))
	q.Set("issuer", issuer)
	expected.RawQuery = q.Encode()
	totp := NewTOTP(secret, 0, int64(period), "", nil)
	got, err := url.Parse(totp.ProvisioningURI(name, issuer))
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
