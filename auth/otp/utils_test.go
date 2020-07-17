package otp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strconv"
	"testing"
)

func TestBuildUri(t *testing.T) {
	issuer := "example"
	name := "test"
	secret, _ := GenerateRandomSecret(16)
	period := 60
	digits := 8
	algorithm := "SHA256"
	expected, _ := url.Parse(fmt.Sprintf("otpauth://%s/%s", Totp, issuer+":"+name))
	q := expected.Query()
	q.Set("secret", secret)
	q.Set("period", strconv.Itoa(period))
	q.Set("digits", strconv.Itoa(digits))
	q.Set("algorithm", algorithm)
	q.Set("issuer", issuer)
	expected.RawQuery = q.Encode()
	got, err := url.Parse(BuildURI(Totp, secret, name, issuer, algorithm, 0, digits, int64(period)))
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	assert.Panics(t, func() {
		BuildURI("", secret, name, issuer, algorithm, 0, digits, int64(period))
	})
}

func TestGenerateRandomSecret(t *testing.T) {
	str, err := GenerateRandomSecret(6)
	assert.Empty(t, str)
	assert.Error(t, err)
}
