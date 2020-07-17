package otp

// HOTP HMAC-based OTP counters
type HOTP struct {
	OTP
	Algorithm string
}

// NewHOTP ...
func NewHOTP(secret string, digits int, algorithm string, algorithmFunc HashFunc) *HOTP {
	return &HOTP{
		OTP:       NewOTP(secret, digits, algorithmFunc),
		Algorithm: algorithm,
	}
}

// At generates the OTP for the given count.
//	count: the OTP HMAC counter
//	returns: OTP string
func (h *HOTP) At(count int) string {
	return h.GenerateOTP(int64(count))
}

// Verify verifies the OTP passed in against the current counter OTP.
//	otp: the OTP to check against
//	counter: the OTP HMAC counter
func (h *HOTP) Verify(otp string, count int) bool {
	return otp == h.At(count)
}

// ProvisioningURI returns the provisioning URI for the OTP.
func (h *HOTP) ProvisioningURI(name, issuer string, initialCount int) string {
	return BuildURI(Hotp,
		h.Secret,
		name,
		issuer,
		h.Algorithm,
		initialCount,
		h.Digits,
		0)
}
