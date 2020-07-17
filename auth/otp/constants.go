package otp

// DefaultDigits OTP digits >= 6 (usually 6-8)
const DefaultDigits = 6

// DefaultInterval default interval for TOTP
const DefaultInterval = 30 // 30 ticks
// ValidWindowThreshold default window threshold for TOTP
const ValidWindowThreshold = 3 * DefaultInterval

type otpType string

// Totp type
const Totp otpType = "totp"

// Hotp type
const Hotp otpType = "hotp"
