package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	for _, l := range levels {
		_ = NewLogger(l)
	}

	_ = NewLogger("test")
}

func TestNewLoggerWithName(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	for _, l := range levels {
		_, err := NewLoggerWithName("test", l)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := NewLoggerWithName("test", "test")
	if err == nil {
		t.Fail()
	}
}
