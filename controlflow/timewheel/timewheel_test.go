package timewheel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	tw := NewTimeWheel("test", DurationMS(time.Second), 100, 10)
	assert.NotNil(t, tw)

	// Add
	// >
	err := tw.AddTimer(1, NewTimer(
		time.Now().Add(time.Minute).UTC(), func(i ...interface{}) {
			fmt.Println("1 called")
		}, []interface{}{}))
	assert.NoError(t, err)
	// < && nil child
	err = tw.AddTimer(2, NewTimer(
		time.Now().Add(time.Millisecond).UTC(), func(i ...interface{}) {
			fmt.Println("2 called")
		}, []interface{}{}))
	assert.NoError(t, err)
	// child
	tw1 := NewTimeWheel("test", DurationMS(time.Millisecond), 100, 10)
	assert.NotNil(t, tw1)
	tw.SetChild(tw1)
	err = tw.AddTimer(3, NewTimer(
		time.Now().Add(time.Millisecond).UTC(), func(i ...interface{}) {
			fmt.Println("3 called")
		}, []interface{}{}))
	assert.NoError(t, err)

	// Remove
	tw.RemoveTimer(1)

	// Run
	tw1.Run()
	time.Sleep(100 * time.Millisecond)
	tw1.Stop()
}
