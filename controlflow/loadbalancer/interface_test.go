package loadbalancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWithConsistentHash(t *testing.T) {
	lb := NewWithConsistentHash()
	_, err := lb.Select([]string{"0.0.0.0", "0.0.0.1", "0.0.0.2"}, "hello")
	assert.NoError(t, err)
}
