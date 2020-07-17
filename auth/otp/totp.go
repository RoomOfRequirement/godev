package otp

import "time"

// TOTP time-based OTP counters
type TOTP struct {
	OTP
	Algorithm string
	Interval  int64 // seconds
}

// NewTOTP ...
func NewTOTP(secret string, digits int, interval int64, algorithm string, algorithmFunc HashFunc) *TOTP {
	if interval < 1 {
		interval = DefaultInterval
	}
	return &TOTP{
		OTP:       NewOTP(secret, digits, algorithmFunc),
		Algorithm: algorithm,
		Interval:  interval,
	}
}

// At generates OTP string of given timestamp and counterOffset
//	timestamp: the time to generate an OTP for
//	counter_offset: the amount of ticks to add to the time counter
//	notice:
//		it's counterOffset (counterOffset * t.Interval),
//		not timeOffset, timeOffset need be processed through function `timeCode()`
func (t *TOTP) At(timestamp, counterOffset int64) string {
	return t.GenerateOTP(t.timeCode(timestamp) + counterOffset)
}

// Now generates OTP string of now
func (t *TOTP) Now() string {
	return t.GenerateOTP(t.timeCode(time.Now().Unix()))
}

// NowWithExpiration generates OTP string of now and return expiration time
func (t *TOTP) NowWithExpiration() (string, time.Time) {
	now := time.Now()
	timeCode := t.timeCode(now.Unix())
	return t.GenerateOTP(timeCode), now.Add(time.Duration(t.Interval) * time.Second)
}

// Verify verifies given otp string/timestamp pair in a given validWindow time counter
//	notice: validWindow is validWindow times of t.Interval (validWindow * t.Interval)
func (t *TOTP) Verify(otp string, timestamp int64, validWindow int64) bool {
	if validWindow != 0 && validWindow < ValidWindowThreshold {
		for i := validWindow * (-1); i < validWindow+1; i++ {
			if otp == t.At(timestamp, i) {
				return true
			}
		}
		return false
	}
	return otp == t.At(timestamp, 0)
}

// ProvisioningURI returns the provisioning URI for the OTP.
//	This can then be encoded in a QR Code and used to provision an OTP app like Google Authenticator.
//	See also: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func (t *TOTP) ProvisioningURI(name, issuer string) string {
	return BuildURI(Totp,
		t.Secret,
		name,
		issuer,
		t.Algorithm,
		0,
		t.Digits,
		t.Interval)
}

// transfer timestamp to timeCode
func (t *TOTP) timeCode(timestamp int64) int64 {
	return timestamp / t.Interval
}
