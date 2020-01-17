package cron

import (
	"errors"
	"go.uber.org/zap"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// New creates a new Cron
func New() Cron {
	return &scheduler{
		jobs:          []*internalJob{},
		jobID:         0,
		newJobChan:    make(chan *internalJob, 1),
		removeJobChan: make(chan int, 1),
		logger:        NewLogger("info"),
		stopOnce:      sync.Once{},
		stopChan:      make(chan struct{}),
		stopped:       1, // need call Start() method to start
		waiter:        sync.WaitGroup{},
	}
}

// internal scheduler
type scheduler struct {
	jobs          []*internalJob
	jobID         int32
	newJobChan    chan *internalJob
	removeJobChan chan int

	logger *zap.Logger

	stopOnce sync.Once
	stopChan chan struct{}
	stopped  int32

	waiter sync.WaitGroup
}

// AddJob adds cron job and returns job id
func (s *scheduler) AddJob(job Job, executedAt time.Time, nextAt func(time.Time) time.Time) (int, error) {
	if s.Stopped() {
		return -1, errors.New("cron has been stopped")
	}
	id := int(atomic.LoadInt32(&s.jobID))
	j := &internalJob{
		id:     id,
		prev:   time.Time{},
		next:   executedAt,
		job:    job,
		nextAt: nextAt,
	}
	s.newJobChan <- j
	atomic.AddInt32(&s.jobID, 1)
	return id, nil
}

// RemoveJob removes one job according to its id
func (s *scheduler) RemoveJob(id int) {
	s.removeJobChan <- id
}

// Start starts the cron
func (s *scheduler) Start() {
	if !s.Stopped() {
		return
	}
	atomic.StoreInt32(&s.stopped, 0)
	go s.run()
}

func (s *scheduler) run() {
	s.logger.Info("Cron Started")
	now := time.Now()

	for {
		// update jobs
		sort.Sort(byTime(s.jobs))
		// check executable job
		var timer *time.Timer
		if len(s.jobs) == 0 {
			// no job -> sleep one day
			timer = time.NewTimer(24 * time.Hour)
		} else if s.jobs[0].next.IsZero() {
			// remove no schedule job
			s.logger.Info("Remove No Schedule Job", zap.Int("id", s.jobs[0].id))
			s.jobs = s.jobs[1:]
			continue
		} else {
			timer = time.NewTimer(s.jobs[0].next.Sub(now))
		}

		select {
		case now = <-timer.C:
			for _, job := range s.jobs {
				// no job available yet
				if job.next.After(now) || job.next.IsZero() {
					break
				}
				s.startJob(job.job)
				job.prev = job.next
				job.next = job.nextAt(now)
				s.logger.Info("Running Job", zap.Int("id", job.id), zap.Time("now", now), zap.Time("next", job.next))
			}
		case newJob := <-s.newJobChan:
			timer.Stop()
			now = time.Now()
			newJob.next = newJob.nextAt(now)
			s.jobs = append(s.jobs, newJob)
			s.logger.Info("Add New Job", zap.Int("id", newJob.id), zap.Time("now", now))
		case <-s.stopChan:
			timer.Stop()
			return
		case id := <-s.removeJobChan:
			timer.Stop()
			now = time.Now()
			s.removeJob(id)
			s.logger.Info("Remove Job", zap.Int("id", id), zap.Time("now", now))
		}
	}
}

func (s *scheduler) startJob(job Job) {
	s.waiter.Add(1)
	go func() {
		defer s.waiter.Done()
		job.Run()
	}()
}

func (s *scheduler) removeJob(id int) {
	jobs := make([]*internalJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		if job.id != id {
			jobs = append(jobs, job)
		}
	}
	s.jobs = jobs
}

// Stop stops the cron
func (s *scheduler) Stop() {
	s.stopOnce.Do(func() {
		atomic.StoreInt32(&s.stopped, 1)
		// wait for running jobs to finish
		go func() {
			s.waiter.Wait()
			close(s.stopChan)
		}()
		<-s.stopChan
		s.logger.Info("Cron Stopped")
	})
}

// Stopped returns true if the cron stopped
func (s *scheduler) Stopped() bool {
	return atomic.LoadInt32(&s.stopped) != 0
}

type internalJob struct {
	id         int
	prev, next time.Time
	job        Job
	nextAt     func(time.Time) time.Time
}

type byTime []*internalJob

// Len to meet sort interface
func (bt byTime) Len() int {
	return len(bt)
}

// Swap to meet sort interface
func (bt byTime) Swap(i, j int) {
	bt[i], bt[j] = bt[j], bt[i]
}

// Less to meet sort interface
func (bt byTime) Less(i, j int) bool {
	// zero means no next, put it to the end
	if bt[i].next.IsZero() {
		return false
	}
	if bt[j].next.IsZero() {
		return true
	}
	return bt[i].next.Before(bt[j].next)
}
