package circuitbreaker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var _ Interface = (*breaker)(nil)

	b := New(NewSettings(0, 0, 0, nil, nil))
	assert.Equal(t, Closed, b.State())

	b.Trip()
	assert.Equal(t, Open, b.State())

	b.Reset()
	assert.Equal(t, Closed, b.State())

	b = New(NewSettings(20, 5*time.Second, 3*time.Second, nil, nil))
	assert.Equal(t, Closed, b.State())
	b.Trip()
	assert.Equal(t, Open, b.State())
	b.Reset()
	assert.Equal(t, Closed, b.State())

	res, err := b.Execute(func() (i interface{}, err error) {
		return 0, nil
	})
	assert.Equal(t, 0, res.(int))
	assert.NoError(t, err)

	// closed -> open
	for i := 0; i < 10; i++ {
		res, err = b.Execute(func() (i interface{}, err error) {
			return 1, errors.New("test")
		})
		assert.Nil(t, res)
		assert.Error(t, err, "test")
	}
	assert.Equal(t, Open, b.State())

	// open -> half-open
	time.Sleep(3 * time.Second)
	assert.Equal(t, HalfOpen, b.State())

	// half-open -> open
	res, err = b.Execute(func() (i interface{}, err error) {
		return 1, errors.New("test")
	})
	assert.Nil(t, res)
	assert.Error(t, err, "test")
	assert.Equal(t, Open, b.State())

	// open -> half-open
	time.Sleep(3 * time.Second)
	assert.Equal(t, HalfOpen, b.State())

	// half-open -> closed
	for i := 0; i < 20; i++ {
		res, err = b.Execute(func() (i interface{}, err error) {
			return 2, nil
		})
		assert.Equal(t, res.(int), 2)
		assert.NoError(t, err)
	}
	assert.Equal(t, Closed, b.State())

	// closed -> closed
	time.Sleep(5 * time.Second)
	assert.Equal(t, Closed, b.State())
}
