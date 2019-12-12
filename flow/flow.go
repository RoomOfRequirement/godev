package flow

// reference: https://rodaine.com/2018/08/x-files-sync-golang/
// code from: https://github.com/rodaine/executor
// thanks for Rodaine's demo!

import (
	"context"
	"errors"
	"golang.org/x/sync/semaphore"
)

// ErrMaxActions ...
var ErrMaxActions = errors.New("exceed allowed maximum actions")

// Flow struct for control flow
type Flow struct {
	MaxActions int64
	// Actions limits concurrent number of actions
	Actions *semaphore.Weighted
	// Calls limits concurrent number of calls
	Calls *semaphore.Weighted
	// base executor
	Executor Executor
}

// NewFlow creates a flow executor and limits it to a maximum concurrent number of calls and actions
func NewFlow(executor Executor, maxCalls, maxActions int64) *Flow {
	return &Flow{
		MaxActions: maxActions,
		Actions:    semaphore.NewWeighted(maxActions),
		Calls:      semaphore.NewWeighted(maxCalls),
		Executor:   executor,
	}
}

// Execute attempts to acquire the semaphores for the concurrent calls and
// actions before delegating to the decorated Executor
// If Execute is called with more actions than maxActions, an error is returned
func (f *Flow) Execute(ctx context.Context, actions ...Action) error {
	actionNum := int64(len(actions))
	if actionNum > f.MaxActions {
		return ErrMaxActions
	}

	if err := f.Calls.Acquire(ctx, 1); err != nil {
		return err
	}
	defer f.Calls.Release(1)

	if err := f.Actions.Acquire(ctx, actionNum); err != nil {
		return err
	}
	defer f.Actions.Release(actionNum)

	return f.Executor.Execute(ctx, actions...)
}
