package batch

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNewTaskN(t *testing.T) {
	tn := NewTaskN(0)
	assert.Equal(t, int32(1), tn.n)
	assert.Equal(t, 0, len(tn.errs))
	assert.Equal(t, 1, cap(tn.errs))

	tn = NewTaskN(5)
	_ = tn.WithContext(context.TODO())
	for i := 0; i < 10; i++ {
		tn.Go(i, func() error {
			return nil
		})
	}
	err := tn.Wait()
	assert.NoError(t, err)

	tn = NewTaskN(5)
	_ = tn.WithContext(context.TODO())
	for i := 0; i < 10; i++ {
		tn.Go(i, func() error {
			return fmt.Errorf("error here")
		})
	}
	err = tn.Wait()
	assert.Error(t, err)

	tn = NewTaskN(5)
	ctx := tn.WithContext(context.TODO())
	for i := 0; i < 10; i++ {
		tn.Go(i, func() error {
			if rand.Intn(10) < 5 {
				return nil
			}
			return fmt.Errorf("error here")
		})
	}
	err = tn.Wait()
	t.Log(err)
	select {
	case <-ctx.Done():
		t.Log("ctx cancelled")
	}
}
