package loadbalancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRR(t *testing.T) {
	rr, err := NewRR(nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidNodes, err)
	assert.Nil(t, rr)

	rr, err = NewRR([]string{"0.0.0.1", "0.0.0.2"})
	assert.NoError(t, err)
	assert.NotNil(t, rr)

	assert.Equal(t, "0.0.0.1", rr.Next())
	assert.Equal(t, "0.0.0.2", rr.Next())
	assert.Equal(t, "0.0.0.1", rr.Next())
	assert.Equal(t, "0.0.0.2", rr.Next())
}
