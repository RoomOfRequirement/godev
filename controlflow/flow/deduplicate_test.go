package flow

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestDeDuplicate(t *testing.T) {
	// this should be test on a multi-core machine or test action should run some time

	exec := DeDuplicate(ConcurrentExecutor{})

	var ct uint32

	addToCt := func(ctx context.Context) error {
		time.Sleep(time.Millisecond) // keep the action in flight
		atomic.AddUint32(&ct, 1)
		return nil
	}

	noop := func(ctx context.Context) error {
		return nil
	}

	err := exec.Execute(context.Background(),
		Named("add", addToCt), // named `add`
		Named("add", addToCt), // named `add`

		Named("add 1", addToCt), // named `add 1`

		ActionFunc(addToCt), // unnamed
		ActionFunc(noop),    // unrelated
	)

	assert.NoError(t, err)
	assert.Equal(t, uint32(3), ct)
}
