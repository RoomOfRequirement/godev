package cron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	var _ Cron = (*scheduler)(nil)

	cron := New()
	assert.True(t, cron.Stopped())
	cron.Start()
	assert.False(t, cron.Stopped())
	cron.Stop()
	assert.True(t, cron.Stopped())
	// add after stop
	id, err := cron.AddJob(TestJob{}, time.Now().Add(1*time.Second), func(t time.Time) time.Time {
		return t.Add(1 * time.Second)
	})
	assert.Equal(t, -1, id)
	assert.Error(t, err)

	// add one job
	cron = New()
	cron.Start()
	cron.Start()
	id, err = cron.AddJob(TestJob{}, time.Now().Add(1*time.Second), func(t time.Time) time.Time {
		return t.Add(1 * time.Second)
	})
	assert.Equal(t, 0, id)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
	cron.RemoveJob(id)
	cron.Stop()
	assert.True(t, cron.Stopped())

	// add multiple jobs
	cron = New()
	cron.Start()
	now := time.Now()
	id, err = cron.AddJob(TestJob{0}, now.Add(1*time.Second), func(t time.Time) time.Time {
		return time.Time{}
	})
	assert.Equal(t, 0, id)
	assert.NoError(t, err)
	id, err = cron.AddJob(TestJob{1}, now.Add(1*time.Second), func(t time.Time) time.Time {
		return time.Time{}
	})
	assert.Equal(t, 1, id)
	assert.NoError(t, err)
	id, err = cron.AddJob(TestJob{2}, now.Add(1*time.Second), func(t time.Time) time.Time {
		return t.Add(1 * time.Second)
	})
	assert.Equal(t, 2, id)
	assert.NoError(t, err)
	id, err = cron.AddJob(TestJob{3}, now.Add(2*time.Second), func(t time.Time) time.Time {
		return t.Add(2 * time.Second)
	})
	assert.Equal(t, 3, id)
	assert.NoError(t, err)
	id, err = cron.AddJob(TestJob{4}, now.Add(1*time.Second), func(t time.Time) time.Time {
		return time.Time{}
	})
	assert.Equal(t, 4, id)
	assert.NoError(t, err)
	cron.RemoveJob(id)
	id, err = cron.AddJob(TestJob{5}, now.Add(3*time.Second), func(t time.Time) time.Time {
		return time.Time{}
	})
	assert.Equal(t, 5, id)
	assert.NoError(t, err)
	time.Sleep(5 * time.Second)
	cron.Stop()
	assert.True(t, cron.Stopped())
}

func TestByTime(t *testing.T) {
	now := time.Now()
	bt := byTime{
		&internalJob{id: 0, next: now.Add(1 * time.Second)},
		&internalJob{id: 1, next: time.Time{}},
		&internalJob{id: 2, next: now.Add(-1 * time.Second)},
		&internalJob{id: 3, next: now.Add(2 * time.Second)},
		&internalJob{id: 4, next: now.Add(5 * time.Second)},
		&internalJob{id: 5, next: now.Add(3 * time.Second)},
	}
	sort.Sort(bt)
	ordered := []int{2, 0, 3, 5, 4, 1}
	for i, b := range bt {
		if b.id != ordered[i] {
			t.Fail()
		}
	}
}

type TestJob struct {
	id int
}

func (tj TestJob) Run() {
	fmt.Println("tj:", tj.id)
}
