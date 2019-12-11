package goroutinepool

// code from: https://github.com/gammazero/workerpool/blob/master/workerpool.go
// http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
// http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

import (
	"goContainer/queue/deque"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// This value is the size of the queue that workers register their
	// availability to the dispatcher.  There may be hundreds of workers, but
	// only a small channel is needed to register some of the workers.
	readyQueueSize = 16

	// If worker pool receives no new work for this period of time, then stop
	// a worker goroutine.
	idleTimeoutSec = 5
)

// Pool struct
type Pool struct {
	maxWorkers int

	timeout time.Duration

	taskChan     chan Task
	readyWorkers chan chan Task

	waitingQ *deque.Deque

	stopOnce sync.Once

	stoppedChan chan struct{}

	queuedTaskNum int32
	stopped       int32
}

// Task type
type Task func()

// New creates a new worker pool with input maxWorkers and a default timeout
func New(maxWorkers int) *Pool {
	return NewPool(maxWorkers, idleTimeoutSec*time.Second)
}

// NewPool creates a new worker pool with input maxWorkers and timeout
func NewPool(maxWorkers int, timeout time.Duration) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	if timeout < idleTimeoutSec*time.Second {
		timeout = idleTimeoutSec * time.Second
	}

	p := &Pool{
		maxWorkers:    maxWorkers,
		timeout:       timeout,
		taskChan:      make(chan Task, 1), // buffered
		readyWorkers:  make(chan chan Task, readyQueueSize),
		waitingQ:      deque.NewDeque(maxWorkers), // auto expand
		stopOnce:      sync.Once{},
		stoppedChan:   make(chan struct{}),
		queuedTaskNum: 0,
		stopped:       0,
	}
	go p.dispatch()
	return p
}

// Submit enqueues a function for a worker to execute
func (p *Pool) Submit(task Task) {
	if task != nil {
		p.taskChan <- task
	}
}

// SubmitWait enqueues the given function and waits for it to be executed
func (p *Pool) SubmitWait(task Task) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskChan <- func() {
		task()
		close(doneChan)
	}
	<-doneChan
}

// Process processes a payload and return the result synchronously
// just pass payload and func as closure
func (p *Pool) Process(payload interface{}, taskFunc func(payload interface{}) interface{}) interface{} {
	resChan := make(chan interface{})
	p.Submit(func() {
		resChan <- taskFunc(payload)
		close(resChan)
	})
	return <-resChan
}

// ProcessAsync processes a payload asynchronously and put the result into resChan
// just pass payload and func as closure
func (p *Pool) ProcessAsync(payload interface{}, taskFunc func(payload interface{}) interface{}, resChan chan<- interface{}) {
	p.Submit(func() {
		resChan <- taskFunc(payload)
		close(resChan)
	})
}

// ProcessFuture simulate the Future idempotent by just simply wrapping un-buffered chan with a function
func (p *Pool) ProcessFuture(payload interface{}, taskFunc func(payload interface{}) interface{}) (future func() interface{}) {
	resChan := make(chan interface{})
	p.Submit(func() {
		resChan <- taskFunc(payload)
		close(resChan)
	})
	return func() interface{} {
		return <-resChan
	}
}

// QueuedTaskNum returns number of waiting tasks
func (p *Pool) QueuedTaskNum() int {
	return int(atomic.LoadInt32(&p.queuedTaskNum))
}

func (p *Pool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(p.timeout)

	var (
		task           Task
		open, wait     bool
		workerTaskChan chan Task
		workerCnt      int
	)
	startReady := make(chan chan Task)

Loop:
	for {
		if !p.waitingQ.Empty() {
			select {
			case task, open = <-p.taskChan:
				// stopped
				if !open {
					break Loop
				}
				if task == nil {
					wait = true
					break Loop
				}

				p.waitingQ.PushBack(task)
			case workerTaskChan = <-p.readyWorkers:
				// ready work request task
				t, _ := p.waitingQ.PopFront()
				workerTaskChan <- t.(Task)
			}
			atomic.StoreInt32(&p.queuedTaskNum, int32(p.waitingQ.Size()))
			continue
		}
		timeout.Reset(p.timeout)
		select {
		case task, open = <-p.taskChan:
			if !open || task == nil {
				break Loop
			}
			// execute task
			select {
			case workerTaskChan = <-p.readyWorkers:
				// dispatch task to ready worker
				workerTaskChan <- task
			default:
				// no ready worker
				// create a new worker if not exceed maxWorkers
				if workerCnt < p.maxWorkers {
					workerCnt++
					go func(t Task) {
						startWorker(startReady, p.readyWorkers)
						// use new worker to execute task when it is start ready
						taskChan := <-startReady
						taskChan <- t
					}(task)
				} else {
					// enqueue task waiting for next ready worker
					p.waitingQ.PushBack(task)
					atomic.StoreInt32(&p.queuedTaskNum, int32(p.waitingQ.Size()))
				}
			}
		case <-timeout.C:
			// reach to worker timeout, kill one ready worker
			if workerCnt > 0 {
				select {
				case workerTaskChan = <-p.readyWorkers:
					// get one ready worker
					close(workerTaskChan)
					workerCnt--
				default:
					// no ready worker, all workers are busy, so do nothing
				}
			}
		}
	}

	// wait for all tasks done
	if wait {
		for p.waitingQ.Size() != 0 {
			// get a ready worker
			workerTaskChan = <-p.readyWorkers
			// give task to it
			t, _ := p.waitingQ.PopFront()
			workerTaskChan <- t.(Task)
			atomic.StoreInt32(&p.queuedTaskNum, int32(p.waitingQ.Size()))
		}
	}

	// stop all remaining workers after they become ready
	for workerCnt > 0 {
		workerTaskChan = <-p.readyWorkers
		close(workerTaskChan)
		workerCnt--
	}
}

func startWorker(startReady, readyWorkers chan chan Task) {
	go func() {
		taskChan := make(chan Task)
		var (
			task Task
			open bool
		)
		// register ready state to start ready chan
		startReady <- taskChan
		for {
			// read task from dispatcher
			task, open = <-taskChan
			if !open {
				break
			}

			// execute task (blocking)
			task()

			// register ready state to ready workers chan
			readyWorkers <- taskChan
		}
	}()
}

// Stop stops the worker pool and waits for only currently running tasks to complete
// Pending tasks that are not currently running are abandoned
// Tasks must not be submitted to the worker pool after calling stop
func (p *Pool) Stop() {
	p.stop(false)
}

// StopUntilAllDone stops the worker pool and waits for all queued tasks to be executed
// No additional tasks may be submitted, but all pending tasks are executed by workers before this function returns
func (p *Pool) StopUntilAllDone() {
	p.stop(true)
}

func (p *Pool) stop(wait bool) {
	p.stopOnce.Do(func() {
		atomic.StoreInt32(&p.stopped, 1)
		if wait {
			// nil as wait flag
			p.taskChan <- nil
		}
		// close taskChan to stop receive new task
		close(p.taskChan)
		// wait for running tasks to finish
		<-p.stoppedChan
	})
}

// Stopped returns true if this worker pool has been stopped
func (p *Pool) Stopped() bool {
	return atomic.LoadInt32(&p.stopped) != 0
}
