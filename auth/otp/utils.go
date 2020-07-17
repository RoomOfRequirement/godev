package otp

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Int64ToBytes ...
//	HMAC counter value to message bytes
func Int64ToBytes(x int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, x)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// HashToCode turns hash bytes to uin32 code
func HashToCode(hashBytes []byte) uint32 {
	// Last byte of hashBytes as index offset
	offset := int64(hashBytes[len(hashBytes)-1] & 0xf)
	// Next 4 bytes starting at the offset as uint32
	b := make([]byte, 4)
	copy(b, hashBytes[offset:offset+4])
	b[0] = b[0] & 0x7f
	return binary.BigEndian.Uint32(b[:])
}

/*
BuildURI returns the provisioning URI for the OTP; works for either TOTP or HOTP.
	This can then be encoded in a QR Code and used to provision the Google Authenticator app.
    For module-internal use.
    See also:
        https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	params:
		otpType: otp type, hotp/totp
		secret: the hotp/totp secret used to generate the URI
    	name: name of the account
    	initial_count: starting counter value used by hotp.
    	issuer: the name of the OTP issuer; this will be the organization title of the OTP entry in Authenticator.
    	algorithm: the algorithm used in the OTP generation.
    	digits: the length of the OTP generated code.
    	period: the number of seconds the OTP generator is set to expire every code, for totp.
    returns: provisioning uri
*/
func BuildURI(otpType otpType, secret, name, issuer, algorithm string, initialCount, digits int, period int64) string {
	if otpType != Totp && otpType != Hotp {
		panic("Unsupported OTP type: " + otpType)
	}
	name = url.QueryEscape(name)
	issuer = url.QueryEscape(issuer)
	uri, _ := url.Parse(fmt.Sprintf("otpauth://%s/%s", otpType, issuer+":"+name))
	q := uri.Query()
	q.Set("secret", secret)
	if otpType == Hotp {
		q.Set("counter", strconv.Itoa(initialCount))
	}
	q.Set("issuer", issuer)
	if algorithm != "" && algorithm != "sha1" {
		q.Set("algorithm", strings.ToUpper(algorithm))
	}
	if digits != 0 && digits != 6 {
		q.Set("digits", strconv.Itoa(digits))
	}
	if period != 0 && period != 30 {
		q.Set("period", strconv.Itoa(int(period)))
	}
	uri.RawQuery = q.Encode()
	return uri.String()
}

// GenerateRandomSecret generates random secret with given bytes length
func GenerateRandomSecret(bytesLen int) (string, error) {
	if bytesLen < 16 {
		return "", errors.New("secret length needs to be at least 16 bytes")
	}
	bs := make([]byte, bytesLen)
	_, err := rand.Read(bs)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bs), nil
}
