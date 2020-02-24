package utils

// Ternary for simulating ternary operator
func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}
