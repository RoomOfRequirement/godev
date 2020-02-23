package timewheel

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNowMS(t *testing.T) {
	now := NowMS()
	time.Sleep(time.Millisecond)
	assert.True(t, almostEqual(now, NowMS()-1, 1))
}

func TestTimeMS(t *testing.T) {
	assert.True(t, almostEqual(NowMS(), TimeMS(time.Now()), 1))
}

func TestDurationMS(t *testing.T) {
	d := time.Millisecond
	assert.Equal(t, d.Milliseconds(), DurationMS(d))
}

func TestNewTimerWithRepeat(t *testing.T) {
	a, b, c := 0, 1, 2
	timer := NewTimerWithRepeat(time.Now().Add(5*time.Millisecond).UTC(),
		10*time.Millisecond,
		0,
		func(i ...interface{}) {
			*i[0].(*int), *i[1].(*int), *i[2].(*int) = 2, 1, 0
		}, []interface{}{&a, &b, &c})
	timer.Run()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 2, a)
	assert.Equal(t, 1, b)
	assert.Equal(t, 0, c)

	cnt := 0
	timer = NewTimerWithRepeat(time.Now().UTC(),
		1*time.Millisecond,
		2,
		func(i ...interface{}) {
			*i[0].(*int)++
		}, []interface{}{&cnt})
	timer.Run()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 2, cnt)
}

func almostEqual(a, b, tolerance int64) bool {
	if a-b < tolerance || a-b > -tolerance {
		return true
	}
	return false
}
