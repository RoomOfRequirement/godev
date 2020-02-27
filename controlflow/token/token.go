package token

import (
	"context"
	"fmt"
	"math/rand"
)

// Token ...
//	works like `drumming flowers (击鼓传花)`,
//	to limit batch of tasks to be done in certain sequence
type Token struct {
	cnt   int
	slots []chan struct{}
}

// New ... at least one slot
func New(cnt int) *Token {
	if cnt < 1 {
		cnt = 1
	}
	t := &Token{
		cnt:   cnt,
		slots: make([]chan struct{}, cnt),
	}
	for i := 0; i < cnt; i++ {
		// actually you can add more token in one slot
		// but the behavior will be more complex when acquire and release
		// realize it when you really need
		t.slots[i] = make(chan struct{}, 1)
		// notice: no pre-fill here
		// t.slots[i] <- struct{}{}
	}
	return t
}

// AcquireFrom acquires token from idx slot
func (t *Token) AcquireFrom(ctx context.Context, idx int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.slots[idx]:
		return nil
	}
}

// Acquire acquires token
func (t *Token) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for _, s := range t.slots {
			select {
			case <-s:
				return nil
			default:
				continue
			}
		}
		return fmt.Errorf("no token here, you should first pass one")
	}
}

// PassTo passes the token to idx slot
//	drumming and pass the flower (token)
func (t *Token) PassTo(ctx context.Context, idx int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.slots[idx] <- struct{}{}:
		return nil
	}
}

// Pass passes the token to any empty slot
//	notice: randomly
func (t *Token) Pass(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.slots[rand.Intn(t.cnt)] <- struct{}{}:
		return nil
	}
}
