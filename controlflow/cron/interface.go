package cron

import "time"

// Job interface
type Job interface {
	Run()
}

// Cron interface
type Cron interface {
	AddJob(job Job, executedAt time.Time, nextAt func(time.Time) time.Time) (int, error)
	RemoveJob(id int)
	Start()
	Stop()
	Stopped() bool
}
