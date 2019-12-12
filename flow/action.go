package flow

import (
	"context"
)

// Action interface
type Action interface {
	Execute(ctx context.Context) error
}

// ActionFunc meets Action interface and makes a function into action
type ActionFunc func(ctx context.Context) error

// Execute function to meet Action / Executor interface
func (af ActionFunc) Execute(ctx context.Context) error {
	return af(ctx)
}

// A NamedAction describes an Action that also has a unique identifier
// This interface is used by the DeDuplicate Executor to prevent duplicate actions
// from running concurrently (like sync.Once)
type NamedAction interface {
	Action

	// ID returns the name for this Action
	// Identical actions should return the same ID value
	ID() string
}

type namedAction struct {
	ActionFunc
	name string
}

func (na namedAction) ID() string {
	return na.name
}

// Named creates a NamedAction from ActionFunc with name
func Named(name string, action ActionFunc) NamedAction {
	return namedAction{
		ActionFunc: action,
		name:       name,
	}
}
