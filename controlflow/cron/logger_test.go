package cron

import "testing"

func TestNewLogger(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	for _, l := range levels {
		_ = NewLogger(l)
	}

	_ = NewLogger("test")
}
