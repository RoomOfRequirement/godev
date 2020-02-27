package token

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	b := New(0)
	assert.Equal(t, 1, b.cnt)
	assert.Equal(t, 1, len(b.slots))
	assert.Equal(t, 1, cap(b.slots[0]))
	assert.Equal(t, 0, len(b.slots[0]))

	b = New(3)
	ctx, cancel := context.WithCancel(context.TODO())
	err := b.Acquire(ctx)
	assert.Error(t, err)
	t.Log(err)
	cancel()
	err = b.Acquire(ctx)
	assert.Error(t, err)
	err = b.AcquireFrom(ctx, 0)
	assert.Error(t, err)
	_ = b.PassTo(ctx, 0)
	err = b.PassTo(ctx, 0)
	t.Log(err)

	b = New(3)
	ctx, cancel = context.WithCancel(context.TODO())
	// first pass one token to certain idx
	err = b.Pass(ctx)
	assert.NoError(t, err)
	err = b.Acquire(ctx)
	assert.NoError(t, err)
	err = b.PassTo(ctx, 0)
	assert.NoError(t, err)
	err = b.AcquireFrom(ctx, 0)
	assert.NoError(t, err)
	err = b.PassTo(ctx, 1)
	assert.NoError(t, err)
	err = b.AcquireFrom(ctx, 1)
	assert.NoError(t, err)

	// full
	_ = b.PassTo(ctx, 0)
	_ = b.PassTo(ctx, 1)
	_ = b.PassTo(ctx, 2)
	cancel()
	err = b.Pass(ctx)
	assert.Error(t, err)
}
