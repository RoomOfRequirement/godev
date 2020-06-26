package server

import "sync"

// ConnectionManager ...
type ConnectionManager struct {
	Counter int
	*sync.WaitGroup
}

// NewConnectionManager ...
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Counter:   0,
		WaitGroup: &sync.WaitGroup{},
	}
}

// Add ...
func (cm *ConnectionManager) Add(delta int) {
	cm.Counter += delta
	cm.WaitGroup.Add(delta)
}

// Done ...
func (cm *ConnectionManager) Done() {
	cm.Counter--
	cm.WaitGroup.Done()
}
