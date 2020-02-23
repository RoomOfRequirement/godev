package utils

import (
	"errors"
	"sync"
	"time"
)

// reference from Twitter

const (
	// TimestampBits ... 41 bits
	// <41, 69.7>, <40, 34.8>, <39, 17.4>
	TimestampBits = 41
	// MachineIDBits ... 10 bits
	MachineIDBits = 10
	// SequenceBits ... 12 bits
	SequenceBits = 12

	// MaxElapsedMS ... max timestamp
	MaxElapsedMS = 1<<TimestampBits - 1
	// MaxMachineID ... max machine id
	MaxMachineID = 1<<MachineIDBits - 1
	// MaxSequenceBits ... max sequence
	MaxSequenceBits = 1<<SequenceBits - 1

	// DefaultStartTime ...
	// 2020-02-02 02:02:02
	DefaultStartTime = 1580608922000
)

// SnowFlakeOptions ...
type SnowFlakeOptions struct {
	// when to start the snowflake
	StartTime time.Time
	// how to get the machine id
	MachineID func() uint16 // uint16 = 2 bytes = 16 bits
}

// SnowFlake ...
type SnowFlake struct {
	sync.Mutex

	startTime int64
	machineID uint16
	sequence  uint16

	elapsedTime int64
}

// NewSnowflake ...
func NewSnowflake(options SnowFlakeOptions) (*SnowFlake, error) {
	// invalid start time
	if options.StartTime.After(time.Now()) {
		return nil, errors.New("invalid start time (after now)")
	}
	if options.MachineID == nil {
		return nil, errors.New("invalid machine id func (nil)")
	}
	if options.MachineID() > MaxMachineID {
		return nil, errors.New("invalid machine id (beyond the max)")
	}
	sf := SnowFlake{
		Mutex:       sync.Mutex{},
		startTime:   0,
		machineID:   options.MachineID(),
		sequence:    0,
		elapsedTime: 0,
	}
	if options.StartTime.IsZero() {
		sf.startTime = DefaultStartTime
	} else {
		sf.startTime = timeToInt64(options.StartTime)
	}
	return &sf, nil
}

// NextUID ...
func (sf *SnowFlake) NextUID() (uint64, error) {
	sf.Lock()
	defer sf.Unlock()
	current := timeToInt64(time.Now()) - sf.startTime
	// first sequence
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else {
		// next sequence number
		sf.sequence = (sf.sequence + 1) & MaxSequenceBits
		// sequence number sold out
		if sf.sequence == 0 {
			sf.elapsedTime++
			// wait for next time
			time.Sleep(time.Millisecond)
		}
	}
	// sf.elapsedTime >= current
	if sf.elapsedTime >= MaxElapsedMS {
		// over time
		return 0, errors.New("timestamp sold out")
	}
	return uint64(sf.elapsedTime<<(MachineIDBits+SequenceBits)) | uint64(sf.machineID<<MachineIDBits) | uint64(sf.sequence), nil
}

func timeToInt64(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}
