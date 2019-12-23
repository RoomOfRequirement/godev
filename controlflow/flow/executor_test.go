package flow

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync/atomic"
	"testing"
)

func TestSequentialExecutor_Execute(t *testing.T) {
	t.Parallel()

	seq := SequentialExecutor{}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)

		n := 10
		actions := make([]Action, n)
		for i := 0; i < n; i++ {
			x := i
			actions[i] = ActionFunc(func(ctx context.Context) error {
				_, _ = fmt.Fprint(buf, x)
				return nil
			})
		}

		err := seq.Execute(context.Background(), actions...)
		assert.NoError(t, err)
		assert.Equal(t, "0123456789", buf.String())
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ct := 0

		addToCt := ActionFunc(func(ctx context.Context) error {
			ct++
			return nil
		})

		actions := []Action{
			addToCt,
			ActionFunc(func(ctx context.Context) error {
				return errors.New("some error")
			}),
			addToCt,
		}

		err := seq.Execute(context.Background(), actions...)
		assert.Error(t, err)
		assert.Equal(t, 1, ct)
	})

	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()

		ct := 0

		addToCt := ActionFunc(func(ctx context.Context) error {
			ct++
			return nil
		})

		actions := []Action{addToCt, addToCt, addToCt}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := seq.Execute(ctx, actions...)
		assert.Equal(t, context.Canceled, err)
		assert.Zero(t, ct)
	})
}

func TestConcurrentExecutor_Execute(t *testing.T) {
	t.Parallel()
	n := runtime.NumCPU()

	exec := ConcurrentExecutor{}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		actionNums := n * 10
		actions := make([]Action, actionNums)
		for i := 0; i < actionNums; i++ {
			actions[i] = addToCt
		}

		err := exec.Execute(context.Background(), actions...)
		assert.NoError(t, err)
		assert.Equal(t, uint32(actionNums), ct)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		waitAct := ActionFunc(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})

		errAct := ActionFunc(func(ctx context.Context) error {
			return errors.New("some error")
		})

		err := exec.Execute(context.Background(), addToCt, waitAct, errAct)
		assert.Error(t, err)
		// addToCt can be executed before or after waitAct
		assert.True(t, 0 == ct || 1 == ct)
	})

	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		actionNums := n * 10
		actions := make([]Action, actionNums)
		for i := 0; i < actionNums; i++ {
			actions[i] = addToCt
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := exec.Execute(ctx, actions...)
		assert.Equal(t, context.Canceled, err)
		assert.Zero(t, ct)
	})
}

func TestPoolExecutor_Execute(t *testing.T) {
	t.Parallel()
	n := runtime.NumCPU()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		exec, done := NewPool(n)
		defer done()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		actionNums := n * 10
		actions := make([]Action, actionNums)
		for i := 0; i < actionNums; i++ {
			actions[i] = addToCt
		}

		err := exec.Execute(context.Background(), actions...)
		assert.NoError(t, err)
		assert.Equal(t, uint32(actionNums), ct)
	})

	t.Run("empty actions", func(t *testing.T) {
		t.Parallel()

		exec, done := NewPool(n)
		defer done()

		err := exec.Execute(context.Background())
		assert.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		t.Parallel()

		exec, done := NewPool(n)
		defer done()

		noopAct := ActionFunc(func(ctx context.Context) error {
			return nil
		})

		waitAct := ActionFunc(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})

		errAct := ActionFunc(func(ctx context.Context) error {
			return errors.New("some error")
		})

		err := exec.Execute(context.Background(), noopAct, waitAct, errAct)
		assert.Error(t, err)
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		exec, done := NewPool(n)
		defer done()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		actionNums := n * 10
		actions := make([]Action, actionNums)
		for i := 0; i < actionNums; i++ {
			actions[i] = addToCt
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := exec.Execute(ctx, actions...)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("done pool", func(t *testing.T) {
		t.Parallel()

		exec, done := NewPool(n)
		done()

		var ct uint32

		addToCt := ActionFunc(func(ctx context.Context) error {
			atomic.AddUint32(&ct, 1)
			return nil
		})

		actionNums := n * 10
		actions := make([]Action, actionNums)
		for i := 0; i < actionNums; i++ {
			actions[i] = addToCt
		}

		err := exec.Execute(context.Background(), actions...)
		assert.Error(t, err)
	})
}
