package roundrobin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRR(t *testing.T) {
	rr := NewSRR()
	err := rr.SetNodes(nil, nil)
	assert.Equal(t, ErrInvalidNodes, err)

	err = rr.SetNodes([]string{"0.0.0.1", "0.0.0.2"}, nil)
	assert.NoError(t, err)

	assert.Equal(t, "0.0.0.1", rr.Next())
	assert.Equal(t, "0.0.0.2", rr.Next())
	assert.Equal(t, "0.0.0.1", rr.Next())
	assert.Equal(t, "0.0.0.2", rr.Next())

	assert.Equal(t, "0.0.0.1", rr.Next())
	rr.Reset()
	assert.Equal(t, "0.0.0.1", rr.Next())
	assert.Equal(t, "0.0.0.2", rr.Next())
}
