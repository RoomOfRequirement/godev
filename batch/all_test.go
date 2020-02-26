package batch

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAll(t *testing.T) {
	a := NewAll(0)
	assert.Equal(t, 1, cap(a.errs))
	a = NewAllWithLimit(0, 0)
	assert.Equal(t, 1, cap(a.errs))
	assert.NotNil(t, a.sema)
	assert.NotNil(t, a.semaCtx)

	a = NewAllWithLimit(5, 2)
	ctx := a.WithContext(context.TODO())
	for i := 0; i < 3; i++ {
		a.Go(i, func() error {
			select {
			case <- ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})
	}
	err := a.Wait()
	assert.NoError(t, err)

	a = NewAllWithLimit(3, 1)
	ctx = a.WithContext(context.TODO())
	for i := 0; i < 3; i++ {
		a.Go(i, func() error {
			select {
			case <- ctx.Done():
				return ctx.Err()
			default:
				return fmt.Errorf("error here")
			}
		})
	}
	err = a.Wait()
	assert.Error(t, err)
	t.Log(err)
}
