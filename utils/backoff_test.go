package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBackOff(t *testing.T) {
	// maxBackOff < initialBackOff -> maxBackOff = initialBackOff
	err := BackOff(100*time.Millisecond, 0, 2, func() error {
		return errors.New("error here")
	})
	assert.Error(t, err)

	// maxCalls = 0, call f only once
	err = BackOff(100*time.Millisecond, time.Second, 0, func() error {
		return errors.New("error here")
	})
	assert.Error(t, err)

	// success
	err = BackOff(100*time.Millisecond, time.Second, 0, func() error {
		return nil
	})
	assert.NoError(t, err)
}
