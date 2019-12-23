package flow

// doc: https://godoc.org/golang.org/x/sync/singleflight

import (
	"context"
	"golang.org/x/sync/singleflight"
)

// DeDuplicate enables a duplicate function call suppression mechanism
// it means the NamedAction with same name will be prevent executing
// at the same time and only prevent at the same time
func DeDuplicate(exec Executor) Executor {
	return deDuplicate{
		exec: exec,
		sf:   new(singleflight.Group),
	}
}

type deDuplicate struct {
	exec Executor
	sf   *singleflight.Group
}

func (dp deDuplicate) Execute(ctx context.Context, actions ...Action) error {
	wrapped := make([]Action, len(actions))

	for i, a := range actions {
		if na, ok := a.(NamedAction); ok {
			wrapped[i] = deDupAction{
				NamedAction: na,
				sf:          dp.sf,
			}
		} else {
			wrapped[i] = a
		}
	}

	return dp.exec.Execute(ctx, wrapped...)
}

type deDupAction struct {
	NamedAction
	sf *singleflight.Group
}

func (dpa deDupAction) Execute(ctx context.Context) error {
	// all actions only return an error, so we don't care if the value is shared or not
	_, err, _ := dpa.sf.Do(dpa.ID(), func() (interface{}, error) {
		return nil, dpa.NamedAction.Execute(ctx)
	})

	return err
}
