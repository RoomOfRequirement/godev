package batch

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErrorN(t *testing.T) {
	en := NewErrorN(0)
	assert.Equal(t, 0, len(en.errs))
	assert.Equal(t, 1, cap(en.errs))

	en.Go(0, func() error {
		return nil
	})
	err := en.Wait()
	assert.NoError(t, err)

	en = NewErrorN(5)
	ctx := en.WithContext(context.TODO())
	for i := 0; i < 10; i++ {
		en.Go(i, func() error {
			return fmt.Errorf("error here")
		})
	}
	err = en.Wait()
	assert.Error(t, err)
	select {
	case <-ctx.Done():
		t.Log("ctx cancelled")
	}
}
